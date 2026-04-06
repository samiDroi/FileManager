import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-confirm-button',
  imports: [],
  templateUrl: './confirm-button-component.html',
  styleUrl: './confirm-button-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ConfirmButtonComponent {
  @Input({ required: true }) text: string = '';
  @Input() icon?: string;
  @Input() isDisabled: boolean = false;
  @Input() isLoading: boolean = false;
  @Input() variant: 'primary' | 'secondary' | 'success' | 'danger' = 'primary'; 
  @Input() type: 'submit' | 'button' = 'submit';
  @Output() clicked = new EventEmitter<void>();

  onButtonClick(): void {
    if (!this.isLoading && !this.isDisabled) {
      //clicked event emitter to parent component
      this.clicked.emit();
    }
  }
}
