import { HttpClient } from '@angular/common/http';
import { inject, Injectable, signal } from '@angular/core';
import { CreateUserRequest, CreatedUserData, UpdateUserRequest, UserData } from '../../../core/domain/dtos/user';
import { ApiRoutes } from '../../../core/routes';
import Swal from '../../../shared/utils/swal';
import { Observable } from 'rxjs';


@Injectable({
  providedIn: 'root',
})
export class UserService {
  private http = inject(HttpClient);
  private usersRequestVersion = 0;
  state = signal({
    users: new Map<number, UserData>(),
    loading: false
  })

  createUser(request: CreateUserRequest): Observable<CreatedUserData> {
    return this.http.post<CreatedUserData>(ApiRoutes.users.create, request);
  }

  getUserDetails(userId: number): Observable<UserData> {
    return this.http.get<UserData>(ApiRoutes.users.getById(userId));
  }

  updateUser(userId: number, request: UpdateUserRequest): Observable<UserData> {
    return this.http.put<UserData>(ApiRoutes.users.update(userId), request);
  }

  upsertUser(user: UserData): void {
    this.usersRequestVersion++;
    this.state.update(current => {
      const users = new Map(current.users);
      users.set(user.id, user);
      return { ...current, users };
    });
  }

  getAllUsers(endpoint = ApiRoutes.users.getAll): void {
    const requestVersion = ++this.usersRequestVersion;
    this.http.get<UserData[]>(endpoint, { withCredentials: true }).subscribe({
      next: (users) => {
        if (requestVersion !== this.usersRequestVersion) {
          return;
        }
        const usersMap = new Map<number, UserData>();
        if (users) {
          users.forEach(user => usersMap.set(user.id, user));
        }
        this.state.set({ users: usersMap, loading: false });
      },
      error: (err) => {
        console.error('Error al obtener usuarios:', err);
        this.state.set({...this.state(), loading: false });
      }
    });
  }

  getAssignableUsers(): void {
    this.getAllUsers(ApiRoutes.users.assignable);
  }

  async deleteUser(userId: number): Promise<void> {
    const result = await Swal.fire({
      title: '¿Eliminar usuario?',
      text: 'Esta acción no se puede deshacer',
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
    });

    if (!result.isConfirmed) return;

    this.http.delete(ApiRoutes.users.delete(userId), { withCredentials: true }).subscribe({
      next: () => {
        const currentState = this.state();
        const updatedUsers = new Map(currentState.users);
        updatedUsers.delete(userId);
        this.state.set({ ...currentState, users: updatedUsers });
        Swal.fire({
          icon: 'success',
          title: 'Usuario eliminado',
          text: 'El usuario se ha eliminado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
      },
      error: (err) => {
        console.error('Error al eliminar usuario:', err);
      }
    });
  }
  searchUsers(query: string): void {
     if (!query.trim()) {
      this.getAllUsers();
      return;
    }
    const requestVersion = ++this.usersRequestVersion;
    this.http.get<UserData[]>(ApiRoutes.users.search(query), { withCredentials: true }).subscribe({
      next: (users) => {
        if (requestVersion !== this.usersRequestVersion) {
          return;
        }
        const usersMap = new Map<number, UserData>();
        if (users) {
          users.forEach(user => usersMap.set(user.id, user));
        }
        this.state.set({ users: usersMap, loading: false });
      },
      error: (err) => {
        console.error('Error al buscar usuarios:', err);
        this.state.set({...this.state(), loading: false });
      }
    });
  }
}
