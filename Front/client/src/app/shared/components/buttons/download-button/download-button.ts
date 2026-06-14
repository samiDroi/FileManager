import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-download-button',
  imports: [],
  templateUrl: './download-button.html',
  styleUrl: './download-button.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DownloadButton {
  @Input() label: string = 'Descargar';
  @Output() clicked = new EventEmitter<void>();

  onClick(): void {
    this.clicked.emit();
  }
}
