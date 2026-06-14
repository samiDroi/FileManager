import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-create-button',
  imports: [],
  templateUrl: './create-button.html',
  styleUrl: './create-button.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class CreateButton {
  @Input() label: string = 'Crear';
  @Output() clicked = new EventEmitter<void>();

  onClick(): void {
    this.clicked.emit();
  }
}
