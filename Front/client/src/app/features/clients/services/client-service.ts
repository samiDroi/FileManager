import { inject, Injectable, signal } from '@angular/core';
import { catchError, EMPTY, map } from 'rxjs';
import { Client } from '../../../core/domain/models/clients';
import { ApiRoutes } from '../../../core/routes/apiUrls';
import { HttpClient } from '@angular/common/http';
import { HttpParams } from '@angular/common/http';
@Injectable({
  providedIn: 'root',
})
export class ClientService {
  private http = inject(HttpClient);

  state = signal({
    clients: new Map<number, Client>(),
    loading: true
  });

  getAllClients(){
    this.http.get<Client[]>(ApiRoutes.clients.getAll,{withCredentials:true}).pipe(
      map(clients => new Map(clients.map(client => [client.id, client]))),
      catchError(error => {
        console.error('Error al obtener clientes:', error);
        this.state.update(state => ({ ...state, loading: false, error: 'No se pudieron cargar los clientes' }));
        return EMPTY; // Corta el flujo para que no llegue al 'next' del subscribe
      })
    ).subscribe({
      // Solo llega aquí si la petición fue exitosa
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

    this.http.get<Client[]>(ApiRoutes.clients.search, { params }).pipe(
      map(clients => new Map(clients.map(client => [client.id, client]))),
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

  createClient(newClient: Omit<Client,'id'>){
    this.state.update(state => ({...state, loading:true}))
    this.http.post<Client>(ApiRoutes.clients.create, newClient).subscribe({
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
      },

      error: (error) => {
        console.error('Error al crear el cliente:', error);
        this.state.update(state => ({ ...state, loading: false }));
      }
    })
  }

  // Método para ACTUALIZAR
  updateClient(id: number, updatedData: Client): void {
    this.state.update(state => ({ ...state, loading: true }));

    // Usamos PUT y le pasamos el ID en la URL
    this.http.put<Client>(ApiRoutes.clients.update(id), updatedData).subscribe({
      next: (updatedClient) => {
        this.state.update(state => {
          const updatedMap = new Map(state.clients);
          
          // Al usar .set() con un ID que ya existe, el Map simplemente 
          // reemplaza los datos viejos con los nuevos. ¡Magia!
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
}
