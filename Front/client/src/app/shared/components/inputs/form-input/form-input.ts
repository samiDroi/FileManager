import { ChangeDetectionStrategy, Component, input, Input } from '@angular/core';
import { FormControl,ReactiveFormsModule } from '@angular/forms';

@Component({
  selector: 'app-form-input',
  imports: [ReactiveFormsModule],
  templateUrl: './form-input.html',
  styleUrl: './form-input.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormInput {
  // Recibimos el control del formulario desde el padre de forma estricta
  @Input({ required: true }) control!: FormControl;
  
  // Opciones de configuración
  @Input() label: string = '';
  @Input() type: 'text' | 'email' | 'password' = 'text';
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
