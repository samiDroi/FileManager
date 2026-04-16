import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { FormComponent } from './form-component/form-component';

@Component({
  selector: 'app-proyects-component',
  imports: [],
  templateUrl: './proyects-component.html',
  styleUrl: './proyects-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ProyectsComponent {
  protected modal = inject(MatDialog)
  
  openCrateModal(){
    const modalRef = this.modal.open(FormComponent,{
      width: "500px",
      data: {proyect:null}
    }

    )
  }
}
