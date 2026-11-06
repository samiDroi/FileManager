import { ChangeDetectionStrategy, Component, inject, signal } from '@angular/core';
import { AppRoutes } from '../../../core/routes/appRoutes';
import { FormBuilder, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { FormInput } from '../../../shared/components/inputs/form-input/form-input';
import { ConfirmButtonComponent } from '../../../shared/components/buttons/confirm-button-component/confirm-button-component';
import { CreateButton } from '../../../shared/components/buttons/create-button/create-button';
import { AuthService } from '../../auth-service';
import { firstValueFrom } from 'rxjs';
@Component({
  selector: 'app-login-component',
  imports: [FormInput,ReactiveFormsModule,ConfirmButtonComponent,CreateButton],
  templateUrl: './login-component.html',
  styleUrl: './login-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class LoginComponent {
  private fb = inject(FormBuilder);
  private router = inject(Router);
  private authService = inject(AuthService);
  errorMessage = signal<string | null>(null);
  // Formulario reactivo
  loginForm = this.fb.group({
    username: ['', [Validators.required]],
    password: ['', [Validators.required]],
  });

  // Getters para pasar los controles de forma segura al componente genérico
  get usernameControl(): FormControl {
    return this.loginForm.get('username') as FormControl;
  }

  get passwordControl(): FormControl {
    return this.loginForm.get('password') as FormControl;
  }

  async onSubmit() {
    // 1. Marcar todo como tocado por si intentan enviar vacío
    if (this.loginForm.invalid) {
      this.loginForm.markAllAsTouched();
      return;
    }

    // 2. Limpiar errores previos al intentar de nuevo
    this.errorMessage.set(null);

    try {
      const response = await firstValueFrom(this.authService.login(this.loginForm.getRawValue()));
      this.loginForm.reset();
      const targetRoute = response?.is_admin
        ? AppRoutes.admin.dashboard
        : AppRoutes.features.dashboard;
      this.router.navigate([targetRoute]);
    } catch (error: any) {
      console.error('Error en login:', error);
      
      // 3. Extraer el mensaje de error del backend (HttpErrorResponse)
      // Si el backend no manda mensaje, ponemos uno por defecto.
      // Regla de seguridad: NUNCA especifiques si falló el usuario o la contraseña, siempre di "Usuario o contraseña incorrectos".
      const msg = error.error?.error || 'Usuario o contraseña incorrectos. Por favor, intenta de nuevo.';
      
      // 4. Actualizar el signal para mostrarlo en el HTML
      this.errorMessage.set(msg);
      
      // 5. Limpiar solo la contraseña
      this.loginForm.patchValue({ password: '' });
    }
  }
  // onSignIn(): void {
  //   if (this.loginForm.invalid) {
  //     this.loginForm.markAllAsTouched();
  //     return;
  //   }

  //   const credentials = this.loginForm.value;
  //   console.log('Datos listos para enviar al backend en Go:', credentials);

  //   this.authService.login(credentials).subscribe({
  //     next: (response) => {
  //       console.log('Login exitoso:', response);
  //       this.router.navigate(['/']);
  //     },
  //     error: (error) => {
  //       console.error('Error en login:', error);
  //     }
  //   });
  // }

  onForgotPassword(): void {
    this.router.navigate([`${AppRoutes.auth.forgotPassword}`]);
  }
//redirect to register page
  onCreateAccount(): void {
    this.router.navigate([`${AppRoutes.auth.register}`]);
  }
}
