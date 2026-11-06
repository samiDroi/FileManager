import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { GroupData, GroupRequest } from '../../../../core/domain/dtos/rbac';
import { CancelButton } from '../../../../shared/components/buttons/cancel-button/cancel-button';
import { ConfirmButtonComponent } from '../../../../shared/components/buttons/confirm-button-component/confirm-button-component';
import { FormInput } from '../../../../shared/components/inputs/form-input/form-input';

interface GroupFormDialogData {
  group: GroupData | null;
}

@Component({
  selector: 'app-group-form-dialog',
  imports: [ReactiveFormsModule, FormInput, CancelButton, ConfirmButtonComponent],
  templateUrl: './group-form-dialog.html',
  styleUrl: './group-form-dialog.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GroupFormDialog {
  private readonly formBuilder = inject(FormBuilder);
  private readonly dialogRef = inject(MatDialogRef<GroupFormDialog>);
  protected readonly data = inject<GroupFormDialogData>(MAT_DIALOG_DATA);

  protected readonly isEdit = Boolean(this.data?.group);
  protected readonly groupForm = this.formBuilder.group({
    name: [this.data?.group?.name ?? '', Validators.required],
    description: [this.data?.group?.description ?? ''],
  });

  protected onCancel(): void {
    this.dialogRef.close(null);
  }

  protected onSubmit(): void {
    if (this.groupForm.invalid) {
      this.groupForm.markAllAsTouched();
      return;
    }

    const request: GroupRequest = {
      name: this.groupForm.get('name')?.value?.trim() ?? '',
      description: this.groupForm.get('description')?.value?.trim() ?? '',
    };
    this.dialogRef.close(request);
  }
}
