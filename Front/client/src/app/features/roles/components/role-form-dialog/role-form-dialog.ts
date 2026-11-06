import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import {
  PermissionData,
  RoleData,
  RoleRequest,
} from '../../../../core/domain/dtos/rbac';
import { CancelButton } from '../../../../shared/components/buttons/cancel-button/cancel-button';
import { ConfirmButtonComponent } from '../../../../shared/components/buttons/confirm-button-component/confirm-button-component';
import { FormInput } from '../../../../shared/components/inputs/form-input/form-input';

export interface RoleFormDialogData {
  role: RoleData | null;
  permissions: PermissionData[];
}

@Component({
  selector: 'app-role-form-dialog',
  imports: [ReactiveFormsModule, FormInput, CancelButton, ConfirmButtonComponent],
  templateUrl: './role-form-dialog.html',
  styleUrl: './role-form-dialog.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RoleFormDialog {
  private readonly formBuilder = inject(FormBuilder);
  private readonly dialogRef = inject(MatDialogRef<RoleFormDialog>);
  protected readonly data = inject<RoleFormDialogData>(MAT_DIALOG_DATA);

  protected readonly isEdit = Boolean(this.data?.role);
  protected readonly selectedPermissionIds = signal<Set<number>>(
    new Set((this.data?.role?.permissions ?? []).map(permission => permission.id))
  );

  protected readonly roleForm = this.formBuilder.group({
    name: [this.data?.role?.name ?? '', Validators.required],
    description: [this.data?.role?.description ?? ''],
  });

  protected isPermissionSelected(permissionId: number): boolean {
    return this.selectedPermissionIds().has(permissionId);
  }

  protected togglePermission(permissionId: number): void {
    this.selectedPermissionIds.update(current => {
      const next = new Set(current);
      if (next.has(permissionId)) {
        next.delete(permissionId);
      } else {
        next.add(permissionId);
      }
      return next;
    });
  }

  protected onCancel(): void {
    this.dialogRef.close(null);
  }

  protected onSubmit(): void {
    if (this.roleForm.invalid) {
      this.roleForm.markAllAsTouched();
      return;
    }

    const request: RoleRequest = {
      name: this.roleForm.get('name')?.value?.trim() ?? '',
      description: this.roleForm.get('description')?.value?.trim() ?? '',
      permission_ids: Array.from(this.selectedPermissionIds()),
    };
    this.dialogRef.close(request);
  }
}
