import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-empty-state',
  standalone: true,
  imports: [],
  templateUrl: './empty-state-component.html',
  styleUrl: './empty-state-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class EmptyStateComponent {
  @Input() titulo: string = '';
  @Input() subtitulo: string = '';
  @Input() textoBoton: string = '';
  @Input() imagenUrl: string = '';
  @Output() buttonClick = new EventEmitter<void>();
}
