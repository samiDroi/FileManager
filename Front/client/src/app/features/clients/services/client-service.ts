import { DestroyRef, inject, Injectable, signal } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { catchError, EMPTY, map, Observable, Subject, switchMap, tap } from 'rxjs';
import { ClientRequest, ClientData, ClientsWithProjects } from '../../../core/domain/dtos';
import { ApiRoutes } from '../../../core/routes/apiUrls';
import { HttpClient, HttpParams } from '@angular/common/http';
import Swal from '../../../shared/utils/swal';

@Injectable({
  providedIn: 'root',
})
export class ClientService {
  private http = inject(HttpClient);
  private destroyRef = inject(DestroyRef);
  private readonly clientListRequests = new Subject<string | null>();

  state = signal({
    clients: new Map<number, ClientData>(),
    loading: true
  });

  stateWithProjects = signal({
    clientsWithProjects: [] as ClientsWithProjects[],
    loading: true
  });

  constructor() {
    this.clientListRequests.pipe(
      tap(() => this.state.update(state => ({ ...state, loading: true }))),
      switchMap(query => this.requestClients(query)),
      takeUntilDestroyed(this.destroyRef)
    ).subscribe(clientsMap => {
      this.state.update(state => ({
        ...state,
        clients: clientsMap,
        loading: false
      }));
    });
  }

  private requestClients(query: string | null): Observable<Map<number, ClientData>> {
    const request$ = query === null
      ? this.http.get<ClientData[]>(ApiRoutes.clients.getAll, { withCredentials: true })
      : this.http.get<ClientData[]>(ApiRoutes.clients.search, {
        params: new HttpParams().set('q', query)
      });

    return request$.pipe(
      map(clients => {
        const clientsMap = new Map<number, ClientData>();
        (clients ?? []).forEach(client => clientsMap.set(client.id, client));
        return clientsMap;
      }),
      catchError(error => {
        console.error('Error al obtener clientes:', error);
        this.state.update(state => ({
          ...state,
          loading: false,
          error: 'No se pudieron cargar los clientes'
        }));
        return EMPTY;
      })
    );
  }

  getAllClients(): void {
    this.clientListRequests.next(null);
  }

  searchClients(query: string): void {
    if (!query.trim()) {
      this.getAllClients();
      return;
    }

    this.clientListRequests.next(query);
  }

  createClient(newClient: ClientRequest){
    this.state.update(state => ({...state, loading:true}))
    this.http.post<ClientData>(ApiRoutes.clients.create, newClient).subscribe({
      next:(createdClient) => {
        this.state.update(state => {
          const updatedMap = new Map(state.clients);
          updatedMap.set(createdClient.id, createdClient);
          return {
            ...state,
            clients: updatedMap,
            loading: false
          };
        });
        this.getClientsWithProjects();
        Swal.fire({
          icon: 'success',
          title: 'Cliente creado',
          text: 'El cliente se ha creado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      },

      error: (error) => {
        console.error('Error al crear el cliente:', error);
        this.state.update(state => ({ ...state, loading: false }));
      }
    })
  }

  updateClient(id: number, updatedData: ClientRequest): void {
    this.state.update(state => ({ ...state, loading: true }));

    this.http.put<ClientData>(ApiRoutes.clients.update(id), updatedData).subscribe({
      next: (updatedClient) => {
        this.state.update(state => {
          const updatedMap = new Map(state.clients);
          updatedMap.set(id, updatedClient);
          return {
            ...state,
            clients: updatedMap,
            loading: false
          };
        });
        Swal.fire({
          icon: 'success',
          title: 'Cliente actualizado',
          text: 'El cliente se ha actualizado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      },
      error: (error) => {
        console.error('Error al actualizar:', error);
        this.state.update(state => ({ ...state, loading: false }));
      }
    });
  }
  getClientsWithProjects(){
    this.http.get<ClientsWithProjects[]>(ApiRoutes.clients.withProjects, {withCredentials:true}).subscribe({
      next: (data) => {
        this.stateWithProjects.update(state => ({
          ...state,
          clientsWithProjects: data,
          loading: false
        }));
      },
      error: (error) => {
        console.error('Error al obtener clientes con proyectos:', error);
        this.stateWithProjects.update(state => ({ ...state, loading: false }));
      }
    });
  }
}
