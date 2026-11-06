import { HttpClient } from '@angular/common/http';
import { computed, inject, Injectable, signal } from '@angular/core';
import { PermissionData, RoleData, RoleRequest } from '../../../core/domain/dtos/rbac';
import { ApiRoutes } from '../../../core/routes';
import Swal from '../../../shared/utils/swal';
import { catchError, Observable, tap, throwError } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class RoleService {
  private readonly http = inject(HttpClient);

  readonly state = signal({
    roles: new Map<number, RoleData>(),
    permissions: [] as PermissionData[],
    loading: false,
  });

  readonly roles = computed(() => Array.from(this.state().roles.values()));
  readonly permissions = computed(() => this.state().permissions);
  readonly loading = computed(() => this.state().loading);

  getAllRoles(): void {
    this.fetchRoles().subscribe({
      error: error => console.error('Error al obtener los roles:', error),
    });
  }

  fetchRoles(): Observable<RoleData[]> {
    this.state.update(current => ({ ...current, loading: true }));
    return this.http.get<RoleData[]>(ApiRoutes.roles.getAll).pipe(
      tap(roles => {
        const rolesMap = new Map<number, RoleData>();
        (roles ?? []).forEach(role => rolesMap.set(role.id, role));
        this.state.update(current => ({ ...current, roles: rolesMap, loading: false }));
      }),
      catchError(error => {
        console.error('Error al obtener los roles:', error);
        this.state.update(current => ({ ...current, loading: false }));
        return throwError(() => error);
      })
    );
  }

  getAllPermissions(): void {
    this.http.get<PermissionData[]>(ApiRoutes.permissions.getAll).subscribe({
      next: permissions => {
        this.state.update(current => ({
          ...current,
          permissions: permissions ?? [],
        }));
      },
      error: error => {
        console.error('Error al obtener los permisos:', error);
      },
    });
  }

  createRole(request: RoleRequest): void {
    this.http.post<RoleData>(ApiRoutes.roles.create, request).subscribe({
      next: createdRole => {
        this.state.update(current => {
          const roles = new Map(current.roles);
          roles.set(createdRole.id, createdRole);
          return { ...current, roles };
        });
        Swal.fire({
          icon: 'success',
          title: 'Rol creado',
          text: 'El rol se ha creado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al crear el rol:', error);
        Swal.fire('Error', 'No se pudo crear el rol', 'error');
      },
    });
  }

  updateRole(roleId: number, request: RoleRequest): void {
    this.http.put<RoleData>(ApiRoutes.roles.update(roleId), request).subscribe({
      next: updatedRole => {
        this.state.update(current => {
          const roles = new Map(current.roles);
          roles.set(updatedRole.id, updatedRole);
          return { ...current, roles };
        });
        Swal.fire({
          icon: 'success',
          title: 'Rol actualizado',
          text: 'El rol se ha actualizado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al actualizar el rol:', error);
        Swal.fire('Error', 'No se pudo actualizar el rol', 'error');
      },
    });
  }

  async deleteRole(role: RoleData): Promise<void> {
    const result = await Swal.fire({
      title: '¿Eliminar rol?',
      text: `Se eliminará el rol "${role.name}".`,
      icon: 'warning',
      showCancelButton: true,
      confirmButtonText: 'Sí, eliminar',
      cancelButtonText: 'Cancelar',
      confirmButtonColor: '#d33',
      cancelButtonColor: '#3085d6',
    });

    if (!result.isConfirmed) {
      return;
    }

    this.http.delete(ApiRoutes.roles.delete(role.id)).subscribe({
      next: () => {
        this.state.update(current => {
          const roles = new Map(current.roles);
          roles.delete(role.id);
          return { ...current, roles };
        });
        Swal.fire({
          icon: 'success',
          title: 'Rol eliminado',
          text: 'El rol se ha eliminado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al eliminar el rol:', error);
        Swal.fire('Error', 'No se pudo eliminar el rol', 'error');
      },
    });
  }
}
