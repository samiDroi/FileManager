import { HttpClient } from '@angular/common/http';
import { Injectable, signal, computed, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { tap } from 'rxjs/operators';
import { User } from '../core/domain/models/users';
import { ApiRoutes } from '../core/routes/apiUrls';

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
  readonly isLoggedIn = computed(() => this.userState() !== null);

  constructor() {
    this.initializeFromStorage();
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
