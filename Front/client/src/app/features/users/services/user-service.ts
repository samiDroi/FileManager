import { HttpClient } from '@angular/common/http';
import { inject, Injectable, signal } from '@angular/core';
import { UserData } from '../../../core/domain/dtos/user';
import { ApiRoutes } from '../../../core/routes';


@Injectable({
  providedIn: 'root',
})
export class UserService {
  private http = inject(HttpClient);
  state = signal({
    users: new Map<number, UserData>(),
    loading: false
  })

  getAllUsers(): void {
    this.http.get<UserData[]>(ApiRoutes.users.getAll, { withCredentials: true }).subscribe({
      next: (users) => {
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

  editUser(user: UserData): void{
    this.http.put<UserData>(ApiRoutes.users.update(user.id), user, { withCredentials: true }).subscribe({
      next:(updatedUser) => {
        const currentState = this.state();
        const updatedUsers = new Map(currentState.users);
        updatedUsers.set(updatedUser.id, updatedUser);
        this.state.set({ ...currentState, users: updatedUsers });
      },
      error: (err) => {
        console.error('Error al editar usuario:', err);
      }
    });
  }

  deleteUser(userId: number): void{
    this.http.delete(ApiRoutes.users.delete(userId), { withCredentials: true }).subscribe({
      next: () => {
        const currentState = this.state();
        const updatedUsers = new Map(currentState.users);
        updatedUsers.delete(userId);
        this.state.set({ ...currentState, users: updatedUsers });
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
    this.http.get<UserData[]>(ApiRoutes.users.search(query), { withCredentials: true }).subscribe({
      next: (users) => {
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
