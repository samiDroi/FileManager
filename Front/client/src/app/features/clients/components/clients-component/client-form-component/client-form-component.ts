import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatDialogModule, MatDialogRef, MAT_DIALOG_DATA } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';

@Component({
  selector: 'app-client-form-component',
  imports: [ReactiveFormsModule],
  templateUrl: './client-form-component.html',
  styleUrl: './client-form-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ClientFormComponent implements OnInit {
  private fb = inject(FormBuilder)
  protected data = inject(MAT_DIALOG_DATA)
  private dialogRef = inject(MatDialogRef<ClientFormComponent>)
  public isEdit = false;

  clientForm = this.fb.group({
    name: [this.data?.client?.name || '', Validators.required],
    email: [this.data?.client?.contact_info?.email || ''],
    contact_phone: [this.data?.client?.contact_info?.contact_phone || '']
  })

  ngOnInit(): void {
    if(this.data && this.data.client){
      this.isEdit = true
      this.clientForm.patchValue({
        name: this.data.client.name,
        email: this.data.client.contact_info?.email || '',
        contact_phone: this.data.client.contact_info?.contact_phone || ''
      })
    }
  }

  onCancel(){
    this.dialogRef.close(null)
  }

  onSubmit(){
    if(this.clientForm.valid){
      const formName = this.clientForm.get('name')?.value || '';
      const formEmail = this.clientForm.get('email')?.value || null;
      const formPhone = this.clientForm.get('contact_phone')?.value || null;
      //payload: important data to send to the database
      const payload: any = { name: formName }

      if (formEmail) payload.email = formEmail;
      if (formPhone) payload.contact_phone = formPhone;

      if (this.isEdit) {
        payload.id = this.data.client.id;
      }

      this.dialogRef.close(payload);
    }
  }
}
