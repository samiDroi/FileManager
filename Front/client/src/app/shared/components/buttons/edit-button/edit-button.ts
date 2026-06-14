import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-edit-button',
  imports: [],
  templateUrl: './edit-button.html',
  styleUrl: './edit-button.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EditButton {
  @Input() label: string = 'Editar';
  @Output() clicked = new EventEmitter<void>();

  onClick(): void {
    this.clicked.emit();
  }
}
