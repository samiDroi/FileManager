import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { ActivatedRoute } from '@angular/router';
import Swal from 'sweetalert2';


@Component({
  selector: 'app-reset-password-component',
  imports: [],
  templateUrl: './reset-password-component.html',
  styleUrl: './reset-password-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ResetPasswordComponent implements OnInit {
  private fb = inject(FormBuilder);
  private routeTkn = inject(ActivatedRoute);
  private token: string | null = null;
  resetForm = this.fb.group({
    password: ['', [Validators.required]],
    newPassword: ['', [Validators.required]],
  });

  ngOnInit(): void {
    this.token = this.routeTkn.snapshot.paramMap.get('token');
  }
   
  onSubmit(): void {
    if (this.resetForm.valid && this.token) {
      const { password, newPassword } = this.resetForm.value;
      console.log('Password reset submitted:', { password, newPassword });
    }
  }
  
}
