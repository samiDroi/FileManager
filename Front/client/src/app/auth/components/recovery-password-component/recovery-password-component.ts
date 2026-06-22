import { ChangeDetectionStrategy, Component, inject, OnInit, signal } from '@angular/core';
import { AbstractControl, FormBuilder, ReactiveFormsModule, ValidationErrors, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import Swal from '../../../shared/utils/swal';
import { AuthService } from '../../auth-service';
import { finalize } from 'rxjs';

@Component({
  selector: 'app-recovery-password-component',
  imports: [ReactiveFormsModule],
  templateUrl: './recovery-password-component.html',
  styleUrl: './recovery-password-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RecoveryPasswordComponent implements OnInit {
  private fb = inject(FormBuilder);
  private routeTkn = inject(ActivatedRoute);
  private token: string | null = null;
  private authService = inject(AuthService);
  private router = inject(Router);

  isLoading = signal(false);
  showPassword = false;
  showConfirm = false;

  recoveryForm = this.fb.group({
    newPassword: ['', [Validators.required, Validators.minLength(6)]],
    confirmPassword: ['', [Validators.required]]
  }, { validators: this.passwordsMatchValidator })

  ngOnInit(): void {
    this.routeTkn.queryParams.subscribe(params => {
      this.token = params['token'];
      if (!this.token) {
        Swal.fire({ icon: 'error', title: 'Token no válido', text: 'El enlace es inválido o expiró.' })
      }
    });
  }

  private passwordsMatchValidator(control: AbstractControl): ValidationErrors | null {
    const password = control.get('newPassword');
    const confirm = control.get('confirmPassword');
    if (password?.value && confirm?.value && password.value !== confirm.value) {
      return { mismatch: true };
    }
    return null;
  }

  togglePassword(): void {
    this.showPassword = !this.showPassword;
  }

  toggleConfirm(): void {
    this.showConfirm = !this.showConfirm;
  }

  onSubmit() {
    if (this.recoveryForm.invalid || !this.token || this.isLoading()) return;

    const newPassword = this.recoveryForm.get('newPassword')?.value;
    if (typeof newPassword !== 'string') return;

    this.isLoading.set(true);

    this.authService.recoveryPassword(this.token, newPassword).pipe(
      finalize(() => this.isLoading.set(false))
    ).subscribe({
      next: () => {
        Swal.fire({ icon: 'success', title: 'Contraseña restablecida', text: 'Tu contraseña ha sido restablecida exitosamente.' })
        this.router.navigate(['/login'])
      },
      error: (err) => {
        Swal.fire({ icon: 'error', title: 'Error', text: err.error?.error || 'Ocurrió un error al restablecer la contraseña.' })
      }
    });
  }
}
