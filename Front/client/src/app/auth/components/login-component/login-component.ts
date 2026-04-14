import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { AppRoutes } from '../../../core/routes/appRoutes';
import { FormBuilder, FormControl, Validators, ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { FormInput } from '../../../shared/components/inputs/form-input/form-input';
import { AuthService } from '../../auth-service';
import { firstValueFrom } from 'rxjs';
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

  async onSubmit(){
    if(this.loginForm.valid){
      try{
        await firstValueFrom(this.authService.login(this.loginForm.getRawValue()));
        this.loginForm.reset();
        this.router.navigate([AppRoutes.features.dashboard]);
      }catch(error){
        console.error('Error en login:', error);
        this.loginForm.patchValue({ password: '' }); // Limpiar solo la contraseña
      }
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
    console.log('Lógica para recuperar contraseña');
  }
//redirect to register page
  onCreateAccount(): void {
    this.router.navigate([`${AppRoutes.auth.register}`]);
  }
}
