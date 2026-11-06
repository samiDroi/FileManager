import { HttpClient } from '@angular/common/http';
import { computed, inject, Injectable, signal } from '@angular/core';
import { catchError, finalize, map, Observable, of, shareReplay } from 'rxjs';
import { UserPermissionsData } from '../domain/dtos/rbac';
import { ApiRoutes } from '../routes';
import { can } from '../authorization/can';

interface PermissionState {
  roles: string[];
  permissions: Set<string>;
  isAdmin: boolean;
  loaded: boolean;
  loading: boolean;
}

@Injectable({
  providedIn: 'root',
})
export class PermissionService {
  private readonly http = inject(HttpClient);
  private readonly state = signal<PermissionState>({
    roles: [],
    permissions: new Set<string>(),
    isAdmin: false,
    loaded: false,
    loading: false,
  });
  private pendingRequest: Observable<boolean> | null = null;

  readonly canManageRoles = computed(() => {
    return this.can('gestionar', 'roles');
  });

  readonly isAdmin = computed(() => this.state().isAdmin);
  readonly roles = computed(() => this.state().roles);
  readonly loaded = computed(() => this.state().loaded);

  can(action: string, resource: string): boolean {
    const current = this.state();
    return can({
      isAdmin: current.isAdmin,
      permissions: Array.from(current.permissions),
    }, action, resource);
  }

  load(force = false): Observable<boolean> {
    if (!force && this.state().loaded) {
      return of(this.canManageRoles());
    }
    if (this.pendingRequest) {
      return this.pendingRequest;
    }

    this.state.update(current => ({ ...current, loading: true }));
    this.pendingRequest = this.http.get<UserPermissionsData>(ApiRoutes.permissions.me).pipe(
      map(response => {
        const permissions = new Set(response?.permissions ?? []);
        this.state.set({
          roles: response?.roles ?? [],
          permissions,
          isAdmin: Boolean(response?.is_admin),
          loaded: true,
          loading: false,
        });
        return this.can('gestionar', 'roles');
      }),
      catchError(error => {
        console.error('Error al obtener los permisos del usuario:', error);
        this.state.update(current => ({ ...current, loaded: true, loading: false }));
        return of(false);
      }),
      finalize(() => {
        this.pendingRequest = null;
      }),
      shareReplay({ bufferSize: 1, refCount: false }),
    );
    return this.pendingRequest;
  }

  setUserPermissions(roles: string[], permissions: string[], isAdmin: boolean): void {
    this.state.set({
      roles,
      permissions: new Set(permissions),
      isAdmin,
      loaded: true,
      loading: false,
    });
  }

  clear(): void {
    this.pendingRequest = null;
    this.state.set({
      roles: [],
      permissions: new Set<string>(),
      isAdmin: false,
      loaded: false,
      loading: false,
    });
  }
}
