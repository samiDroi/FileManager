import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { FormBuilder, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialog, MatDialogRef } from '@angular/material/dialog';
import { Proyect } from '../../../../core';

@Component({
  selector: 'app-form-component',
  imports: [],
  templateUrl: './form-component.html',
  styleUrl: './form-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormComponent implements OnInit{
  private fb = inject(FormBuilder)
  protected data = inject(MAT_DIALOG_DATA);
  //reference of the dialog/modal opened 
  private modalRef = inject(MatDialogRef<Proyect>)
  public isEdit = false;

  proyectForm = this.fb.group({
    //patch stored data or empty
    name: [this.data?.proyect?.name || '',Validators.required]
  })
  ngOnInit(): void {
    if(this.data || this.data?.proyect){
      this.isEdit = true
      this.proyectForm.patchValue({
        name: this.data.proyect.name
      })
    }
  }

  onSubmit(): void{
    if(!this.proyectForm.valid){
      const nameForm = this.proyectForm.get("name")?.value || '';
      
      const payload: any = {name: nameForm}
      if(this.isEdit){
        //if isEdit is true, match the form data with the current proyect selected
        payload.id = this.data.proyect.id;
      }
      //close modal and charge data
      this.modalRef.close(payload);
    }
  }
  onCancel(): void{
    this.modalRef.close(null)
  }
}
