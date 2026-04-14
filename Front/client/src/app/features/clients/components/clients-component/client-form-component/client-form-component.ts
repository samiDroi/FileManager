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
  //im using container presentational but the formulary is so simple so i injected this here
  private fb = inject(FormBuilder)
  protected data = inject(MAT_DIALOG_DATA)
  private dialogRef = inject(MatDialogRef<ClientFormComponent>)
  public isEdit = false;

  clientForm = this.fb.group({
    name: [this.data?.client?.name || '',Validators.required]
  })

  ngOnInit(): void {
    if(this.data && this.data.client){
      this.isEdit = true

      this.clientForm.patchValue({
        name: this.data.client.name
      })
    }
  }

  onCancel(){
    this.dialogRef.close(null)
  }

onSubmit(){
    if(this.clientForm.valid){
      
      // 1. Extraemos el valor asegurándonos de que es un string
      const formName = this.clientForm.get('name')?.value || '';

      // 2. Si es MODO EDICIÓN, mandamos el ID y el Name
      if (this.isEdit) {
        const updatePayload = {
          id: this.data.client.id,
          name: formName
        };
        this.dialogRef.close(updatePayload);
      } 
      // 3. Si es MODO CREACIÓN, mandamos SOLO el Name (Justo lo que Go pide)
      else {
        const createPayload = {
          name: formName
        };
        this.dialogRef.close(createPayload);
      }
    }
  }
}
