import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-cancel-button',
  imports: [],
  templateUrl: './cancel-button.html',
  styleUrl: './cancel-button.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CancelButton {
  @Input() label: string = 'Cancelar';
  @Output() clicked = new EventEmitter<void>();

  onClick(): void {
    this.clicked.emit();
  }
}
