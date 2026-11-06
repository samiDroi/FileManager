import { HttpClient } from '@angular/common/http';
import { computed, inject, Injectable, signal } from '@angular/core';
import { Observable, tap } from 'rxjs';
import {
  AddGroupFileRequest,
  AddGroupUserRequest,
  AvailableUserData,
  GroupData,
  GroupFileData,
  GroupRequest,
  GroupUserData,
} from '../../../core/domain/dtos/rbac';
import { ApiRoutes } from '../../../core/routes';
import Swal from '../../../shared/utils/swal';

@Injectable({
  providedIn: 'root',
})
export class GroupService {
  private readonly http = inject(HttpClient);

  readonly state = signal({
    groups: new Map<number, GroupData>(),
    loading: false,
  });

  readonly groups = computed(() => Array.from(this.state().groups.values()));
  readonly loading = computed(() => this.state().loading);

  getAllGroups(): void {
    this.state.update(current => ({ ...current, loading: true }));
    this.http.get<GroupData[]>(ApiRoutes.groups.getAll).subscribe({
      next: groups => {
        const groupsMap = new Map<number, GroupData>();
        (groups ?? []).forEach(group => groupsMap.set(group.id, group));
        this.state.set({ groups: groupsMap, loading: false });
      },
      error: error => {
        console.error('Error al obtener los grupos:', error);
        this.state.update(current => ({ ...current, loading: false }));
      },
    });
  }

  createGroup(request: GroupRequest): void {
    this.http.post<GroupData>(ApiRoutes.groups.create, request).subscribe({
      next: createdGroup => {
        this.state.update(current => {
          const groups = new Map(current.groups);
          groups.set(createdGroup.id, createdGroup);
          return { ...current, groups };
        });
        Swal.fire({
          icon: 'success',
          title: 'Grupo creado',
          text: 'El grupo se ha creado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al crear el grupo:', error);
        Swal.fire('Error', 'No se pudo crear el grupo', 'error');
      },
    });
  }

  updateGroup(groupId: number, request: GroupRequest): void {
    this.http.put<GroupData>(ApiRoutes.groups.update(groupId), request).subscribe({
      next: updatedGroup => {
        this.state.update(current => {
          const groups = new Map(current.groups);
          groups.set(updatedGroup.id, updatedGroup);
          return { ...current, groups };
        });
        Swal.fire({
          icon: 'success',
          title: 'Grupo actualizado',
          text: 'El grupo se ha actualizado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al actualizar el grupo:', error);
        Swal.fire('Error', 'No se pudo actualizar el grupo', 'error');
      },
    });
  }

  async deleteGroup(group: GroupData): Promise<void> {
    const result = await Swal.fire({
      title: '¿Eliminar grupo?',
      text: `Se eliminará el grupo "${group.name}" y sus accesos asociados.`,
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

    this.http.delete(ApiRoutes.groups.delete(group.id)).subscribe({
      next: () => {
        this.state.update(current => {
          const groups = new Map(current.groups);
          groups.delete(group.id);
          return { ...current, groups };
        });
        Swal.fire({
          icon: 'success',
          title: 'Grupo eliminado',
          text: 'El grupo se ha eliminado correctamente',
          timer: 2000,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al eliminar el grupo:', error);
        Swal.fire('Error', 'No se pudo eliminar el grupo', 'error');
      },
    });
  }

  getGroupUsers(groupId: number): Observable<GroupUserData[]> {
    return this.http.get<GroupUserData[]>(ApiRoutes.groups.getUsers(groupId));
  }

  getAvailableUsers(groupId: number): Observable<AvailableUserData[]> {
    return this.http.get<AvailableUserData[]>(ApiRoutes.groups.getAvailableUsers(groupId));
  }

  addGroupUser(groupId: number, request: AddGroupUserRequest): Observable<GroupUserData> {
    return this.http
      .post<GroupUserData>(ApiRoutes.groups.addUser(groupId), request)
      .pipe(tap(() => this.adjustGroupCounts(groupId, 0, 1)));
  }

  removeGroupUser(groupId: number, userId: number): Observable<void> {
    return this.http
      .delete<void>(ApiRoutes.groups.removeUser(groupId, userId))
      .pipe(tap(() => this.adjustGroupCounts(groupId, 0, -1)));
  }

  getGroupFiles(groupId: number): Observable<GroupFileData[]> {
    return this.http.get<GroupFileData[]>(ApiRoutes.groups.getFiles(groupId));
  }

  addGroupFile(groupId: number, request: AddGroupFileRequest): Observable<void> {
    return this.http
      .post<void>(ApiRoutes.groups.addFile(groupId), request)
      .pipe(tap(() => this.adjustGroupCounts(groupId, 1, 0)));
  }

  revokeGroupFile(groupId: number, fileId: number): Observable<void> {
    return this.http
      .delete<void>(ApiRoutes.groups.revokeFile(groupId, fileId))
      .pipe(tap(() => this.adjustGroupCounts(groupId, -1, 0)));
  }

  adjustGroupCounts(groupId: number, fileDelta: number, userDelta: number): void {
    this.state.update(current => {
      const group = current.groups.get(groupId);
      if (!group) {
        return current;
      }

      const groups = new Map(current.groups);
      groups.set(groupId, {
        ...group,
        file_count: Math.max(0, group.file_count + fileDelta),
        user_count: Math.max(0, group.user_count + userDelta),
      });
      return { ...current, groups };
    });
  }
}
