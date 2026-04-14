import { HttpClient } from '@angular/common/http';
import { computed, inject, Injectable, signal } from '@angular/core';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { User } from '../core/domain/models/users';
import { LocalManagerService, LocalKeys } from '../core/services/local-manager-service';
import { ApiRoutes } from '../core/routes/apiUrls';

interface AuthResponse {
  id: number;
  username: string;
  email: string;
  token: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private http = inject(HttpClient);
  private localManager = inject(LocalManagerService);
  
  // 🌟 1. LA ÚNICA FUENTE DE VERDAD (Privada)
  private userState = signal<User | null>(null);

  // 🌟 2. LO QUE EXPONEMOS AL MUNDO (Solo lectura)
  public readonly currentUser = this.userState.asReadonly();

  // 🌟 3. ESTADO DERIVADO (Computed)
  public readonly isLoggedIn = computed(() => this.userState() !== null);

  constructor() {
    // Inicializar con el token guardado si existe
    this.initializeFromStoredToken();
  }

  // Inicializar el estado si hay un token guardado
  private initializeFromStoredToken(): void {
    const token = this.localManager.getElement(LocalKeys.token);
    // Leemos el usuario guardado como texto
    const storedUserString = this.localManager.getElement(LocalKeys.userData);

    if (token && storedUserString) {
      try {
        // Convertimos el texto de vuelta a un objeto User real
        const parsedUser = JSON.parse(storedUserString) as User;
        this.userState.set(parsedUser);
      } catch (e) {
        // Si por alguna razón el texto se corrompió, cerramos la sesión
        this.logout();
      }
    }
  }

  register(data: any): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(ApiRoutes.auth.register, data).pipe(
      tap((response) => {
        this.setToken(response);
      })
    );
  }

  login(credentials: any): Observable<AuthResponse> {
    return this.http.post<AuthResponse>(ApiRoutes.auth.login, credentials).pipe(
      tap((response) => {
        this.setToken(response);
      })
    );
  }

  logout(): void {
    // Limpiar localStorage completo (token y userData)
    this.localManager.clearStorage();
    
    // Limpiar la fuente de verdad
    this.userState.set(null);
  }

  // Obtener el token actual
  getToken(): string | null {
    return this.localManager.getElement(LocalKeys.token);
  }

  setToken(response: AuthResponse): void {
    // 1. Guardar el token en localStorage
    this.localManager.setElement(LocalKeys.token, response.token);
    
    // 2. Armar el objeto User
    const userDataToSave: User = {
      id: response.id,
      username: response.username,
      email: response.email,
      createdBy: 0,
      updatedBy: 0
    };

    // 3. Convertirlo a texto y guardarlo (usando tu enum LocalKeys)
    const stringUsuario = JSON.stringify(userDataToSave);
    this.localManager.setElement(LocalKeys.userData, stringUsuario);
    
    // 4. Actualizar el Signal en vivo (¡y solo lo hacemos una vez!)
    this.userState.set(userDataToSave);
  }
}