import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { ProjectData, Client } from '../../../../core';
import { ClientService } from '../../../clients/services';
import { CancelButton } from '../../../../shared/components/buttons/cancel-button/cancel-button';
import { ConfirmButtonComponent } from '../../../../shared/components/buttons/confirm-button-component/confirm-button-component';

@Component({
  selector: 'app-form-component',
  imports: [ReactiveFormsModule, CancelButton, ConfirmButtonComponent],
  templateUrl: './form-component.html',
  styleUrl: './form-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class FormComponent implements OnInit{
  private fb = inject(FormBuilder)
  private clientService = inject(ClientService)
  protected data = inject(MAT_DIALOG_DATA);
  private modalRef = inject(MatDialogRef<ProjectData>)
  public isEdit = false;

  clients = computed(() => Array.from(this.clientService.state().clients.values()));

  projectForm = this.fb.group({
    name: [this.data?.project?.name || '', Validators.required],
    client_id: [this.data?.project?.client_id || null, Validators.required]
  })

  ngOnInit(): void {
    this.clientService.getAllClients();
    
    if(this.data?.project){
      this.isEdit = true;
      this.projectForm.patchValue({
        name: this.data.project.name,
        client_id: this.data.project.client_id
      });
    }
  }

  onSubmit(): void {
    if(this.projectForm.valid){
      const nameForm = this.projectForm.get('name')?.value || '';
      const clientIdForm = this.projectForm.get('client_id')?.value;

      const payload: any = {
        name: nameForm,
        client_id: clientIdForm
      };

      if(this.isEdit && this.data?.project){
        payload.id = this.data.project.id;
      }

      this.modalRef.close(payload);
    }
  }

  onCancel(): void{
    this.modalRef.close(null);
  }
}
