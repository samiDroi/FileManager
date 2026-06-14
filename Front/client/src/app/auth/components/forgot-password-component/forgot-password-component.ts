import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { AuthService } from '../../auth-service';
import Swal from 'sweetalert2';
import { finalize } from 'rxjs';

@Component({
  selector: 'app-forgot-password',
  imports: [ReactiveFormsModule],
  templateUrl: './forgot-password-component.html',
  styleUrl: './forgot-password-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ForgotPasswordComponent {
  private authService = inject(AuthService);
  private fb = inject(FormBuilder);

  isLoading = signal(false);

  forgotForm = this.fb.group({
    email: ['', [Validators.required, Validators.email]]
  });

  get emailControl() {
    return this.forgotForm.get('email')!;
  }

  onSubmit(): void {
    if (!this.forgotForm.valid || this.isLoading()) return;

    const email = this.forgotForm.get('email')?.value;
    if (!email) return;

    this.isLoading.set(true);

    this.authService.sendResetPasswordEmail(email).pipe(
      finalize(() => this.isLoading.set(false))
    ).subscribe({
      next: () => {
        Swal.fire({
          icon: 'success',
          title: '¡Correo enviado!',
          text: 'Si el correo existe en nuestro sistema, recibirás un email con instrucciones para restablecer tu contraseña. Si no lo recibes en unos minutos, revisa tu carpeta de spam.',
        })
      },
      error: () => {
        Swal.fire({
          icon: 'error',
          title: 'Error',
          text: 'Ocurrió un error al enviar el email. Por favor, intenta nuevamente más tarde.',
        });
      }
    });
  }
}
