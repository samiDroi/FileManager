import { ChangeDetectionStrategy, Component, inject, OnInit, signal } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { UserData } from '../../../../core/domain/dtos/user';
import { PermissionService } from '../../../../core/services/permission-service';
import { RoleService } from '../../../roles/services/role-service';
import { UserService } from '../../services/user-service';

@Component({
  selector: 'app-edit-user-dialog',
  imports: [ReactiveFormsModule],
  templateUrl: './edit-user-dialog.html',
  styleUrl: './edit-user-dialog.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EditUserDialog implements OnInit {
  private readonly fb = inject(NonNullableFormBuilder);
  protected readonly dialogRef = inject(MatDialogRef<EditUserDialog>);
  protected readonly user = inject<UserData>(MAT_DIALOG_DATA);
  private readonly roleService = inject(RoleService);
  private readonly userService = inject(UserService);
  private readonly permissionService = inject(PermissionService);

  readonly roles = this.roleService.roles;
  readonly loadingRoles = this.roleService.loading;
  readonly canManageRoles = this.permissionService.canManageRoles;
  readonly rolesError = signal(false);
  readonly loadingUser = signal(true);
  readonly submitting = signal(false);
  readonly submitError = signal('');
  readonly currentRoles = signal<UserData['roles']>([]);

  readonly editForm = this.fb.group({
    username: [this.user.username, [Validators.required]],
    email: [this.user.email, [Validators.required, Validators.email]],
    role_id: [0],
  });

  ngOnInit(): void {
    this.userService.getUserDetails(this.user.id).subscribe({
      next: user => {
        this.currentRoles.set(user.roles ?? []);
        if (user.roles?.length) {
          this.editForm.controls.role_id.setValue(user.roles[0].id);
        }
        this.loadingUser.set(false);
      },
      error: error => {
        console.error('No se pudieron cargar los datos del usuario:', error);
        this.loadingUser.set(false);
        this.rolesError.set(true);
      },
    });

    if (this.canManageRoles()) {
      this.roleService.fetchRoles().subscribe({
        next: roles => {
          this.rolesError.set(roles.length === 0);
          if (roles.length > 0 && !this.currentRoles()?.length) {
            this.editForm.controls.role_id.setValue(roles[0].id);
          }
        },
        error: error => {
          console.error('No se pudieron cargar los roles para editar el usuario:', error);
          this.rolesError.set(true);
        },
      });
    }
  }

  onSubmit(): void {
    if (this.editForm.invalid || this.loadingUser() || this.submitting()) {
      this.editForm.markAllAsTouched();
      return;
    }

    const { username, email, role_id } = this.editForm.getRawValue();
    const request = this.canManageRoles()
      ? { username, email, role_id: role_id > 0 ? role_id : undefined }
      : { username, email };

    this.submitting.set(true);
    this.submitError.set('');
    this.userService.updateUser(this.user.id, request).subscribe({
      next: updatedUser => {
        this.dialogRef.close(updatedUser);
      },
      error: error => {
        console.error('Error al actualizar el usuario:', error);
        this.submitting.set(false);
        this.submitError.set(error.error?.message ?? 'No se pudo actualizar el usuario. Verifica los datos e inténtalo de nuevo.');
      },
    });
  }
}
