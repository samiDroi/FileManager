import { inject, Injectable, signal } from '@angular/core';
import { catchError, EMPTY, map } from 'rxjs';
import { ClientRequest, ClientData, ClientsWithProjects } from '../../../core/domain/dtos';
import { ApiRoutes } from '../../../core/routes/apiUrls';
import { HttpClient, HttpParams } from '@angular/common/http';

@Injectable({
  providedIn: 'root',
})
export class ClientService {
  private http = inject(HttpClient);

  state = signal({
    clients: new Map<number, ClientData>(),
    loading: true
  });

  stateWithProjects = signal({
    clientsWithProjects: [] as ClientsWithProjects[],
    loading: true
  });

  getAllClients(){
    this.http.get<ClientData[]>(ApiRoutes.clients.getAll,{withCredentials:true}).pipe(
      map(clients => {
        const clientsMap = new Map<number, ClientData>();
        if (clients) {
          clients.forEach(client => clientsMap.set(client.id, client));
        }
        return clientsMap;
      }),
      catchError(error => {
        console.error('Error al obtener clientes:', error);
        this.state.update(state => ({ ...state, loading: false, error: 'No se pudieron cargar los clientes' }));
        return EMPTY;
      })
    ).subscribe({
      next: (clientsMap) => {
        this.state.update(state => ({ 
          ...state, 
          clients: clientsMap, 
          loading: false 
        }));
      }
    });
  }

  searchClients(query: string): void {
    if (!query.trim()) {
      this.getAllClients();
      return;
    }

    const params = new HttpParams().set('q', query);

    this.http.get<ClientData[]>(ApiRoutes.clients.search, { params }).pipe(
      map(clients => {
        const clientsMap = new Map<number, ClientData>();
        if (clients) {
          clients.forEach(client => clientsMap.set(client.id, client));
        }
        return clientsMap;
      }),
      catchError(error => {
        console.error('Error al buscar clientes:', error);
        return EMPTY;
      })
    ).subscribe({
      next: (clientsMap) => {
        this.state.update(state => ({ 
          ...state, 
          clients: clientsMap
        }));
      }
    });
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