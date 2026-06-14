import { HttpClient } from '@angular/common/http';
import { Injectable, signal, computed, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { catchError, map, tap } from 'rxjs/operators';
import { User } from '../core/domain/models/users';
import { ApiRoutes } from '../core/routes/apiUrls';
import { Observable, of } from 'rxjs';

interface AuthResponse {
  id: number;
  username: string;
  email: string;
}

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private platformId = inject(PLATFORM_ID);
  private http = inject(HttpClient);
  
  private userState = signal<User | null>(null);
  private initialized = signal(false);

  readonly currentUser = this.userState.asReadonly();
  //if the user is logged, true, if not, false
  // readonly isLoggedIn = computed(() => this.userState() !== null);
  //Restore the user session from local storage, to manage the reload pages of the user and the session expiration
  constructor() {
    //without this, the user will need to login to any refresh of the page
    this.initializeFromStorage();
  }
 async verifySession(): Promise<void> {
  if (!this.hasStoredUser()) {
    return; // No hay datos que verificar
  }

  try {
    const res = await fetch(ApiRoutes.auth.refreshToken, {
      method: 'POST',
      credentials: 'include', // envía cookies
    });
    if (!res.ok) {
      this.logout(); // 401 → limpia todo
    }
    // Si es OK (200), no haces nada porque initializeFromStorage() ya seteó userState
  } catch {
    this.logout(); // Error de red (server apagado) → limpia todo
  }
}
  private initializeFromStorage(): void {
    if (!isPlatformBrowser(this.platformId)) {
      return;
    }
    
    const storedUser = localStorage.getItem('userData');
    if (storedUser) {
      try {
        const parsedUser = JSON.parse(storedUser) as User;
        this.userState.set(parsedUser);
      } catch {
        localStorage.removeItem('userData');
      }
    }
    this.initialized.set(true);
  }

  private saveUser(response: AuthResponse): void {
    const userData: User = {
      id: response.id,
      username: response.username,
      email: response.email,
      createdBy: 0,
      updatedBy: 0
    };
    if (isPlatformBrowser(this.platformId)) {
      localStorage.setItem('userData', JSON.stringify(userData));
    }
    this.userState.set(userData);
  }

  login(credentials: any) {
    return this.http.post<AuthResponse>(ApiRoutes.auth.login, credentials).pipe(
      tap(response => this.saveUser(response))
    );
  }

  register(data: any) {
    return this.http.post<AuthResponse>(ApiRoutes.auth.register, data).pipe(
      tap(response => this.saveUser(response))
    );
  }

  logout(): void {
    if (isPlatformBrowser(this.platformId)) {
      localStorage.removeItem('userData');
    }
    this.userState.set(null);
  }

  hasStoredUser(): boolean {
    if (!isPlatformBrowser(this.platformId)) {
      return false;
    }
    return localStorage.getItem('userData') !== null;
  }
}
