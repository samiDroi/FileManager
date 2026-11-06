import { ChangeDetectionStrategy, Component, inject, OnInit, signal } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { RoleService } from '../../../roles/services/role-service';
import { UserService } from '../../services/user-service';

@Component({
  selector: 'app-create-user-dialog',
  imports: [ReactiveFormsModule],
  templateUrl: './create-user-dialog.html',
  styleUrl: './create-user-dialog.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CreateUserDialog implements OnInit {
  private readonly fb = inject(NonNullableFormBuilder);
  protected readonly dialogRef = inject(MatDialogRef<CreateUserDialog>);
  private readonly roleService = inject(RoleService);
  private readonly userService = inject(UserService);

  readonly roles = this.roleService.roles;
  readonly loadingRoles = this.roleService.loading;
  readonly rolesError = signal(false);
  readonly submitting = signal(false);
  readonly submitError = signal('');

  readonly registerForm = this.fb.group({
    username: ['', [Validators.required]],
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]],
    role_id: [0, [Validators.required, Validators.min(1)]],
  });

  ngOnInit(): void {
    this.roleService.fetchRoles().subscribe({
      next: roles => {
        this.rolesError.set(roles.length === 0);
        if (roles.length > 0) {
          this.registerForm.controls.role_id.setValue(roles[0].id);
        }
      },
      error: error => {
        console.error('No se pudieron cargar los roles para el alta de usuario:', error);
        this.rolesError.set(true);
      },
    });
  }

  onSubmit(): void {
    if (this.registerForm.invalid || this.submitting()) {
      this.registerForm.markAllAsTouched();
      return;
    }

    this.submitting.set(true);
    this.submitError.set('');
    this.userService.createUser(this.registerForm.getRawValue()).subscribe({
      next: createdUser => {
        const role = createdUser.roles[0];
        if (!role) {
          console.error('La respuesta de alta no incluye el rol asignado.', createdUser);
          this.submitting.set(false);
          this.dialogRef.close(createdUser);
          return;
        }

        this.dialogRef.close(createdUser);
      },
      error: error => {
        console.error('Error al crear el usuario:', error);
        this.submitting.set(false);
        this.submitError.set(error.error?.message ?? 'No se pudo crear el usuario. Verifica los datos e inténtalo de nuevo.');
      },
    });
  }
}
