import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { AppRoutes } from '../../../core/constants.routes';
import { FormBuilder, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { FormInput } from '../../../shared/components/inputs/form-input/form-input';
import { AuthService } from '../../auth-service';
@Component({
  selector: 'app-login-component',
  imports: [FormInput,ReactiveFormsModule],
  templateUrl: './login-component.html',
  styleUrl: './login-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private authService = inject(AuthService);

  // Formulario reactivo
  signInForm = this.fb.group({
    email: ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required]],
  });

  // Getters para pasar los controles de forma segura al componente genérico
  get emailControl(): FormControl {
    return this.signInForm.get('email') as FormControl;
  }

  get passwordControl(): FormControl {
    return this.signInForm.get('password') as FormControl;
  }

  onSignIn(): void {
    if (this.signInForm.invalid) {
      this.signInForm.markAllAsTouched();
      return;
    }

    const credentials = this.signInForm.value;
    console.log('Datos listos para enviar al backend en Go:', credentials);

    this.authService.login(credentials).subscribe({
      next: (response) => {
        console.log('Login exitoso:', response);
        this.router.navigate(['/']);
      },
      error: (error) => {
        console.error('Error en login:', error);
      }
    });
  }

  onForgotPassword(): void {
    console.log('Lógica para recuperar contraseña');
  }
//redirect to register page
  onCreateAccount(): void {
    this.router.navigate([`${AppRoutes.auth.register}`]);
  }
}
