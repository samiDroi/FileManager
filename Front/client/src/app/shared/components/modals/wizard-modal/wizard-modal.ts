import { ChangeDetectionStrategy, Component, computed, inject, signal, TemplateRef, Type } from '@angular/core';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

@Component({
  selector: 'app-wizard-modal',
  imports: [],
  templateUrl: './wizard-modal.html',
  styleUrl: './wizard-modal.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class WizardModal {
  private modalRef = inject(MatDialogRef<any>);
  //steps is the templates array, mat dialog data is the data to the modal
  private data: { steps: Type<any>[] } = inject(MAT_DIALOG_DATA);
  
  // Estado
  currentIndex = signal(0);
  
  // Computados reactivos
  currentComponent = computed(() => this.data.steps[this.currentIndex()]);
  hasPrevious = computed(() => this.currentIndex() > 0);
  hasNext = computed(() => this.currentIndex() < this.data.steps.length - 1);

  next(): void {
    if (this.hasNext()) {
      this.currentIndex.update(i => i + 1);
    }
  }

  previous(): void {
    if (this.hasPrevious()) {
      this.currentIndex.update(i => i - 1);
    }
  }

  onClose(): void {
    this.modalRef.close();
  }
}
