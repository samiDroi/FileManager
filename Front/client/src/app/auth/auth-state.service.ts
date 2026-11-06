import { HttpClient } from '@angular/common/http';
import { Injectable, signal, computed, PLATFORM_ID, inject } from '@angular/core';
import { isPlatformBrowser } from '@angular/common';
import { finalize, map, shareReplay, tap } from 'rxjs/operators';
import { User } from '../core/domain/models/users';
import { ApiRoutes } from '../core/routes/apiUrls';
import { PermissionService } from '../core/services/permission-service';
import { firstValueFrom, Observable } from 'rxjs';

interface AuthResponse {
  id: number;
  username: string;
  email: string;
  access_token: string;
  is_admin: boolean;
  roles: string[];
  permissions: string[];
}

interface RefreshResponse {
  access_token: string;
  is_admin: boolean;
  roles: string[];
  permissions: string[];
}

@Injectable({
  providedIn: 'root'
})
export class AuthStateService {
  private platformId = inject(PLATFORM_ID);
  private http = inject(HttpClient);
  private permissionService = inject(PermissionService);
  
  private userState = signal<User | null>(null);
  private initialized = signal(false);
  private accessToken = signal<string | null>(null);
  private refreshRequest?: Observable<string>;

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
    await firstValueFrom(this.refreshAccessToken());
  } catch {
    this.logout();
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
    this.permissionService.clear();
    this.userState.set(userData);
    this.accessToken.set(response.access_token);
    this.permissionService.setUserPermissions(
      response.roles ?? [],
      response.permissions ?? [],
      response.is_admin ?? false
    );
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
    this.permissionService.clear();
    this.userState.set(null);
    this.accessToken.set(null);
  }

  getAccessToken(): string | null {
    return this.accessToken();
  }

  refreshAccessToken(): Observable<string> {
    if (!this.refreshRequest) {
      this.refreshRequest = this.http.post<RefreshResponse>(
        ApiRoutes.auth.refreshToken,
        {},
        { withCredentials: true }
      ).pipe(
        tap(response => {
          this.accessToken.set(response.access_token);
          this.permissionService.setUserPermissions(
            response.roles ?? [],
            response.permissions ?? [],
            response.is_admin ?? false
          );
        }),
        map(response => response.access_token),
        finalize(() => this.refreshRequest = undefined),
        shareReplay({ bufferSize: 1, refCount: false })
      );
    }
    return this.refreshRequest;
  }

  hasStoredUser(): boolean {
    if (!isPlatformBrowser(this.platformId)) {
      return false;
    }
    return localStorage.getItem('userData') !== null;
  }
}
