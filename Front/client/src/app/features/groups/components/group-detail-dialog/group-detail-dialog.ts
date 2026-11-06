import { ChangeDetectionStrategy, Component, computed, inject, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatTabsModule } from '@angular/material/tabs';
import {
  AvailableUserData,
  GroupData,
  GroupFileData,
  GroupUserData,
} from '../../../../core/domain/dtos/rbac';
import { GroupService } from '../../services';
import { FileService } from '../../../files/services';
import { PermissionService } from '../../../../core/services/permission-service';
import Swal from '../../../../shared/utils/swal';

export interface GroupDetailDialogData {
  group: GroupData;
}

@Component({
  selector: 'app-group-detail-dialog',
  imports: [MatTabsModule, FormsModule],
  templateUrl: './group-detail-dialog.html',
  styleUrl: './group-detail-dialog.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GroupDetailDialog implements OnInit {
  private readonly groupService = inject(GroupService);
  private readonly fileService = inject(FileService);
  private readonly permissionService = inject(PermissionService);
  private readonly dialogRef = inject(MatDialogRef<GroupDetailDialog>);
  protected readonly data = inject<GroupDetailDialogData>(MAT_DIALOG_DATA);

  protected readonly group = signal<GroupData>(this.data.group);
  protected readonly users = signal<GroupUserData[]>([]);
  protected readonly files = signal<GroupFileData[]>([]);
  protected readonly availableUsers = signal<AvailableUserData[]>([]);
  protected readonly loadingUsers = signal(true);
  protected readonly loadingFiles = signal(true);
  protected readonly selectedUserId = signal<number | null>(null);
  protected readonly selectedRoleInGroup = signal('lector');
  protected readonly selectedFileId = signal<number | null>(null);
  protected readonly selectedFileAccess = signal('leer');
  protected readonly rolesInGroup = [
    { value: 'lector', label: 'Lector' },
    { value: 'editor', label: 'Editor' },
    { value: 'admin', label: 'Administrador del grupo' },
  ];
  protected readonly accessTypes = ['leer', 'descargar', 'editar'];

  protected readonly canManage = computed(() => this.group().can_manage === true);
  protected readonly canEditFiles = computed(
    () => this.canManage() || ['admin', 'editor'].includes(this.group().role_in_group ?? '')
  );
  protected readonly canManageAccess = computed(() => this.permissionService.can('otorgar', 'accesos'));

  protected readonly availableFiles = computed(() => {
    const assigned = new Set(this.files().map(file => file.id));
    return this.fileService.allFiles().filter(file => !assigned.has(file.id));
  });

  protected readonly canAddUser = computed(() => this.selectedUserId() !== null);

  ngOnInit(): void {
    this.loadUsers();
    this.loadFiles();
    if (this.canManage()) {
      this.loadAvailableUsers();
    }
    if (this.canEditFiles() && this.canManageAccess()) {
      this.fileService.onSearch('');
    }
  }

  private loadUsers(): void {
    this.loadingUsers.set(true);
    this.groupService.getGroupUsers(this.group().id).subscribe({
      next: users => {
        this.users.set(users ?? []);
        this.loadingUsers.set(false);
      },
      error: error => {
        console.error('Error al obtener los usuarios del grupo:', error);
        this.loadingUsers.set(false);
      },
    });
  }

  private loadAvailableUsers(): void {
    this.groupService.getAvailableUsers(this.group().id).subscribe({
      next: users => this.availableUsers.set(users ?? []),
      error: error => console.error('Error al obtener los usuarios disponibles:', error),
    });
  }

  private loadFiles(): void {
    this.loadingFiles.set(true);
    this.groupService.getGroupFiles(this.group().id).subscribe({
      next: files => {
        this.files.set(files ?? []);
        this.loadingFiles.set(false);
      },
      error: error => {
        console.error('Error al obtener los archivos del grupo:', error);
        this.loadingFiles.set(false);
      },
    });
  }

  protected onUserSelected(userId: number | null): void {
    this.selectedUserId.set(userId);
  }

  protected onRoleSelected(role: string): void {
    this.selectedRoleInGroup.set(role);
  }

  protected onFileSelected(fileId: number | null): void {
    this.selectedFileId.set(fileId);
  }

  protected onFileAccessSelected(accessType: string): void {
    this.selectedFileAccess.set(accessType);
  }

  protected onAddUser(): void {
    const userId = this.selectedUserId();
    if (userId === null) {
      return;
    }

    this.groupService
      .addGroupUser(this.group().id, {
        user_id: userId,
        role_in_group: this.selectedRoleInGroup(),
      })
      .subscribe({
        next: user => {
          this.users.update(current => [...current, user]);
          this.availableUsers.update(current => current.filter(item => item.id !== user.id));
          this.selectedUserId.set(null);
          this.selectedRoleInGroup.set('lector');
        },
        error: error => {
          console.error('Error al agregar el usuario al grupo:', error);
          Swal.fire('Error', 'No se pudo agregar el usuario al grupo', 'error');
        },
      });
  }

  protected onRemoveUser(user: GroupUserData): void {
    this.groupService.removeGroupUser(this.group().id, user.id).subscribe({
      next: () => {
        this.users.update(current => current.filter(item => item.id !== user.id));
        this.availableUsers.update(current => [...current, {
          id: user.id,
          username: user.username,
          email: user.email,
        }]);
      },
      error: error => {
        console.error('Error al remover el usuario del grupo:', error);
        Swal.fire('Error', 'No se pudo remover el usuario del grupo', 'error');
      },
    });
  }

  protected onAddFile(): void {
    const fileId = this.selectedFileId();
    if (fileId === null) {
      return;
    }

    this.groupService
      .addGroupFile(this.group().id, {
        file_id: fileId,
        access_type: this.selectedFileAccess(),
      })
      .subscribe({
        next: () => {
          this.selectedFileId.set(null);
          this.selectedFileAccess.set('leer');
          this.loadFiles();
          Swal.fire({
            icon: 'success',
            title: 'Archivo agregado',
            timer: 1500,
            showConfirmButton: false,
          });
        },
        error: error => {
          console.error('Error al agregar el archivo al grupo:', error);
          Swal.fire('Error', 'No se pudo agregar el archivo al grupo', 'error');
        },
      });
  }

  protected onRevokeFile(file: GroupFileData): void {
    this.groupService.revokeGroupFile(this.group().id, file.id).subscribe({
      next: () => {
        this.files.update(current => current.filter(item => item.id !== file.id));
        Swal.fire({
          icon: 'success',
          title: 'Acceso revocado',
          text: `El grupo ya no tiene acceso a "${file.tittle}".`,
          timer: 1800,
          showConfirmButton: false,
        });
      },
      error: error => {
        console.error('Error al revocar el acceso al archivo:', error);
        Swal.fire('Error', 'No se pudo revocar el acceso al archivo', 'error');
      },
    });
  }

  protected close(): void {
    this.dialogRef.close();
  }

  protected roleLabel(role: string): string {
    return this.rolesInGroup.find(option => option.value === role)?.label ?? role;
  }
}
