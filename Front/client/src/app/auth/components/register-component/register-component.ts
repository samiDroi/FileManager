import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { NonNullableFormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { MatDialogClose } from '@angular/material/dialog';
import { ConfirmButtonComponent, FormInput, } from '../../../shared/components';
import { AuthService } from '../../auth-service';
import { AppRoutes } from '../../../core/routes/appRoutes';
import Swal from '../../../shared/utils/swal';


@Component({
  selector: 'app-register-component',
  imports: [FormInput,ReactiveFormsModule,ConfirmButtonComponent,MatDialogClose],
  templateUrl: './register-component.html',
  styleUrl: './register-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})

export class RegisterComponent {
  private fb = inject(NonNullableFormBuilder);
  private router = inject(Router);
  private authService = inject(AuthService);

  registerForm = this.fb.group({
    username: ['', [Validators.required]],
    email:    ['', [Validators.required, Validators.email]],
    password: ['', [Validators.required, Validators.minLength(6)]]
  })

  onSubmit(): void {
    if (this.registerForm.invalid) {
      this.registerForm.markAllAsTouched();
      return;
    }
    
    const formData = this.registerForm.getRawValue();
    
    this.authService.register(formData).subscribe({
      next: (response) => {
        Swal.fire({
          icon: 'success',
          title: 'Usuario creado',
          text: 'El usuario se ha registrado correctamente',
          timer: 2000,
          showConfirmButton: false
        });
        this.router.navigate([AppRoutes.auth.login]);
      },
      error: (error) => {
        console.error('Error en registro:', error);
      }
    });
  }
}
