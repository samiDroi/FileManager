import { ChangeDetectionStrategy, Component, computed, inject, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import {
  FileAccessData,
  FileAccessGroupData,
  FileAccessUserData,
  FileData,
} from '../../../../core';
import { FileService } from '../../services';
import { UserService } from '../../../users/services/user-service';
import { GroupService } from '../../../groups/services';
import Swal from '../../../../shared/utils/swal';

export interface FileAccessDialogData {
  file: FileData;
  afterUpload?: boolean;
}

type AccessTab = 'users' | 'groups';

@Component({
  selector: 'app-file-access-dialog',
  imports: [FormsModule],
  templateUrl: './file-access-dialog.html',
  styleUrl: './file-access-dialog.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FileAccessDialog implements OnInit {
  private readonly fileService = inject(FileService);
  private readonly userService = inject(UserService);
  private readonly groupService = inject(GroupService);
  private readonly dialogRef = inject(MatDialogRef<FileAccessDialog>);
  protected readonly data = inject<FileAccessDialogData>(MAT_DIALOG_DATA);

  protected readonly access = signal<FileAccessData>({ users: [], groups: [] });
  protected readonly loading = signal(true);
  protected readonly activeTab = signal<AccessTab>('users');
  protected readonly selectedUserId = signal<number | null>(null);
  protected readonly selectedGroupId = signal<number | null>(null);
  protected readonly selectedUserAccess = signal('leer');
  protected readonly selectedGroupAccess = signal('leer');
  protected readonly accessTypes = ['leer', 'descargar', 'editar'];

  protected readonly availableUsers = computed(() => {
    const granted = new Set(this.access().users.map(entry => entry.user_id));
    return Array.from(this.userService.state().users.values()).filter(user => !granted.has(user.id));
  });

  protected readonly availableGroups = computed(() => {
    const granted = new Set(this.access().groups.map(entry => entry.group_id));
    return this.groupService.groups().filter(group => !granted.has(group.id));
  });

  ngOnInit(): void {
    this.userService.getAssignableUsers();
    this.groupService.getAllGroups();
    this.loadAccess();
  }

  protected loadAccess(): void {
    this.loading.set(true);
    this.fileService.getFileAccess(this.data.file.id).subscribe({
      next: access => {
        this.access.set({
          users: access?.users ?? [],
          groups: access?.groups ?? [],
        });
        this.loading.set(false);
      },
      error: error => {
        console.error('Error al obtener los accesos del archivo:', error);
        this.loading.set(false);
      },
    });
  }

  protected onGrantUser(): void {
    const userId = this.selectedUserId();
    if (userId === null) {
      return;
    }

    this.fileService
      .grantUserAccess(this.data.file.id, {
        user_id: userId,
        access_type: this.selectedUserAccess(),
      })
      .subscribe({
        next: () => {
          this.selectedUserId.set(null);
          this.loadAccess();
          Swal.fire({
            icon: 'success',
            title: 'Acceso otorgado',
            timer: 1500,
            showConfirmButton: false,
          });
        },
        error: error => {
          console.error('Error al otorgar acceso al usuario:', error);
          Swal.fire('Error', 'No se pudo otorgar el acceso al usuario', 'error');
        },
      });
  }

  protected onGrantGroup(): void {
    const groupId = this.selectedGroupId();
    if (groupId === null) {
      return;
    }

    this.fileService
      .grantGroupAccess(this.data.file.id, {
        group_id: groupId,
        access_type: this.selectedGroupAccess(),
      })
      .subscribe({
        next: () => {
          this.selectedGroupId.set(null);
          this.loadAccess();
          Swal.fire({
            icon: 'success',
            title: 'Acceso otorgado',
            timer: 1500,
            showConfirmButton: false,
          });
        },
        error: error => {
          console.error('Error al otorgar acceso al grupo:', error);
          Swal.fire('Error', 'No se pudo otorgar el acceso al grupo', 'error');
        },
      });
  }

  protected onRevokeUser(entry: FileAccessUserData): void {
    this.fileService.revokeUserAccess(this.data.file.id, entry.user_id).subscribe({
      next: () => this.loadAccess(),
      error: error => {
        console.error('Error al revocar el acceso:', error);
        Swal.fire('Error', 'No se pudo revocar el acceso', 'error');
      },
    });
  }

  protected onRevokeGroup(entry: FileAccessGroupData): void {
    this.fileService.revokeGroupAccess(this.data.file.id, entry.group_id).subscribe({
      next: () => this.loadAccess(),
      error: error => {
        console.error('Error al revocar el acceso:', error);
        Swal.fire('Error', 'No se pudo revocar el acceso', 'error');
      },
    });
  }

  protected close(): void {
    this.dialogRef.close();
  }
}
