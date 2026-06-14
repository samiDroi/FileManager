import { CommonModule } from '@angular/common';
import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { App } from '../../../app';
import { AppRoutes } from '../../../core/routes/appRoutes';
import { CancelButton } from '../../../shared/components/buttons/cancel-button/cancel-button';

@Component({
  selector: 'app-sidebar',
  imports: [RouterModule, CancelButton],
  templateUrl: './sidebar-component.html',
  styleUrl: './sidebar-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {
  protected readonly routes = AppRoutes.features;
  // Recibe la info del padre
  @Input() userData: any = null; 

  // Avisa al padre cuando el usuario quiere salir
  @Output() onLogout = new EventEmitter<void>();

  emitLogout() {
    this.onLogout.emit();
  }
}
