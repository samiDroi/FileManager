import { ChangeDetectionStrategy, Component, EventEmitter, input, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AppRoutes } from '../../../core/routes/appRoutes';
import { ClientData, ClientsWithProjects } from '../../../core';
import { CreateButton } from '../../../shared/components/buttons/create-button/create-button';

@Component({
  selector: 'app-sidebar',
  imports: [RouterModule, FormsModule,CreateButton],
  templateUrl: './sidebar-component.html',
  styleUrl: './sidebar-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {
  protected readonly routes = AppRoutes.features;
  clientsList = input<ClientData[] | null>(null);
  @Input() userData: any = null;

  @Output() onLogout = new EventEmitter<void>();
  @Output() onCreateClient = new EventEmitter<void>();

  searchTerm = '';

  get filteredClients(): ClientData[] {
    if (!this.searchTerm.trim()) {
      return this.clientsList() ?? [];
    }
    const term = this.searchTerm.toLowerCase();
    return (this.clientsList() ?? []).filter(c =>
      c.name.toLowerCase().includes(term)
    );
  }

  emitLogout() {
    this.onLogout.emit();
  }

  emitCreateClient() {
    this.onCreateClient.emit();
  }
}
