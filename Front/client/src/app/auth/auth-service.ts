import { HttpClient } from '@angular/common/http';
import { computed, inject, Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { User } from '../core/domain/models/users';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private http = inject(HttpClient);
  private apiUrl = `${environment.apiUrl}/auth`;

  // 🌟 1. LA ÚNICA FUENTE DE VERDAD (Privada)
  // Usamos un Signal. Inicialmente es null (no hay sesión).
  // Es privado para que ningún componente pueda hacer authService.userState.set(...) por accidente.
  private userState = signal<User | null>(null);

  // 🌟 2. LO QUE EXPONEMOS AL MUNDO (Solo lectura)
  // Los componentes pueden leer esto, pero no modificarlo.
  public readonly currentUser = this.userState.asReadonly();

  // 🌟 3. ESTADO DERIVADO (Computed)
  // Angular calcula esto automáticamente. Si currentUser cambia, isLoggedIn se actualiza solo.
  public readonly isLoggedIn = computed(() => this.userState() !== null);
  // public readonly isSuperAdmin = computed(() => this.userState()?.role === 'superadmin');

  // --- MÉTODOS HTTP ---

  register(data: any): Observable<User> {
    return this.http.post<User>(`${this.apiUrl}/register`, data).pipe(
      // 'tap' de RxJS nos permite espiar la respuesta exitosa ANTES de que llegue al componente
      tap((userResponse) => {
        // Actualizamos nuestra fuente de verdad con los datos que mandó Go
        this.userState.set(userResponse);
      })
    );
  }

  login(credentials: any): Observable<User> {
    return this.http.post<User>(`${this.apiUrl}/login`, credentials).pipe(
      tap((userResponse) => {
        this.userState.set(userResponse);
      })
    );
  }

  logout(): void {
    // Limpiamos la fuente de verdad. Automáticamente isLoggedIn pasará a false.
    this.userState.set(null);
  }
}
