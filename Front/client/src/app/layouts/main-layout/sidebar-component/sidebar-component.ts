import { ChangeDetectionStrategy, Component, computed, EventEmitter, inject, input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AppRoutes } from '../../../core/routes/appRoutes';
import { ClientData, ClientsWithProjects } from '../../../core';
import { CreateButton } from '../../../shared/components/buttons/create-button/create-button';
import { PermissionService } from '../../../core/services/permission-service';

@Component({
  selector: 'app-sidebar',
  imports: [RouterModule, FormsModule,CreateButton],
  templateUrl: './sidebar-component.html',
  styleUrl: './sidebar-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class SidebarComponent {
  private readonly permissionService = inject(PermissionService);
  protected readonly routes = AppRoutes.features;
  protected readonly canCreateClients = computed(() => this.permissionService.can('crear', 'clientes'));
  clientsList = input<ClientData[] | null>(null);
  collapsed = input(false);
  hasProjects = input(false);
  projectsLoaded = input(false);

  @Output() onCreateClient = new EventEmitter<void>();
  @Output() requestExpand = new EventEmitter<void>();

  searchTerm = '';

  handleNavigation(): void {
    if (this.collapsed()) {
      this.requestExpand.emit();
    }
  }

  get filteredClients(): ClientData[] {
    if (!this.searchTerm.trim()) {
      return this.clientsList() ?? [];
    }
    const term = this.searchTerm.toLowerCase();
    return (this.clientsList() ?? []).filter(c =>
      c.name.toLowerCase().includes(term)
    );
  }

  emitCreateClient() {
    this.onCreateClient.emit();
  }
}
