import { ChangeDetectionStrategy, Component, inject, output, signal } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, FormControl, ReactiveFormsModule } from '@angular/forms';
import { FormInput } from '../../../../../shared/components/inputs/form-input/form-input';
import { CreateButton } from '../../../../../shared/components/buttons/create-button/create-button';
import { DeleteButton } from '../../../../../shared/components/buttons/delete-button/delete-button';
import { CancelButton } from '../../../../../shared/components/buttons/cancel-button/cancel-button';
import { ConfirmButtonComponent } from '../../../../../shared/components/buttons/confirm-button-component/confirm-button-component';

type ItemFormGroup = FormGroup<{
  id: FormControl<number | null>;
  name: FormControl<string | null>;
}>;

@Component({
  selector: 'app-types-form',
  imports: [ReactiveFormsModule, FormInput, CreateButton, DeleteButton, CancelButton, ConfirmButtonComponent],
  templateUrl: './form-component.html',
  styleUrl: './form-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormComponent {
  private fb = inject(FormBuilder);
  onSave = output<{name: string}[]>();
  onCancel = output<void>();

  form = this.fb.group({
    items: this.fb.array<ItemFormGroup>([])
  });

  items = signal(this.form.controls.items.controls);

  private currentId = 1;

  addItem() {
    const newGroup = this.fb.group({
      id: [this.currentId++],
      name: ['']
    });

    this.form.controls.items.push(newGroup as ItemFormGroup);
    this.items.set([...this.form.controls.items.controls]);
  }

  removeItem(index: number) {
    this.form.controls.items.removeAt(index);
    this.items.set([...this.form.controls.items.controls]);
  }

  onSubmit(): void {
    if (this.form.valid) {
      const rawValues = this.form.value.items || [];
      
      const dataToEmit = rawValues
        .filter(item => item.name)
        .map(item => ({ name: item.name || '' }));

      if (dataToEmit.length > 0) {
        this.onSave.emit(dataToEmit);
      }
    }
  }
}