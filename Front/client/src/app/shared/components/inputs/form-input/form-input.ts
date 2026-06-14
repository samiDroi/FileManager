import { ChangeDetectionStrategy, Component, Input } from '@angular/core';
import { AbstractControl, FormControl, ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-form-input',
  imports: [ReactiveFormsModule],
  templateUrl: './form-input.html',
  styleUrl: './form-input.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormInput {
  @Input({ required: true }) set control(value: AbstractControl) {
    this._control = value as FormControl;
  }
  _control!: FormControl;
  
  // Opciones de configuración
  @Input() label: string = '';
  @Input() type: 'text' | 'email' | 'password' | 'number' = 'text';
  @Input() placeholder: string = '';
  @Input() id: string = Math.random().toString(36).substring(2, 9); // ID único auto-generado

  isPasswordVisible = false;

  togglePasswordVisibility(): void {
    this.isPasswordVisible = !this.isPasswordVisible;
  }

  // Helper para saber qué tipo de input renderizar
  get inputType(): string {
    if (this.type === 'password') {
      return this.isPasswordVisible ? 'text' : 'password';
    }
    return this.type;
  }
}
