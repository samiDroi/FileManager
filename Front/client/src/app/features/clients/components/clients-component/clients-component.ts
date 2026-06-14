import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { ClientService } from '../../services/client-service';
import { ListClientsComponent } from './list-clients-component/list-clients-component';
import { UiSearchBarComponent } from '../../../../shared/components/search-bar/ui-search-bar-component';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
import { MatDialog } from '@angular/material/dialog';
import { ClientFormComponent } from './client-form-component/client-form-component';
import { FormComponent } from '../../../project/projects-component/form-component/form-component';
import { ClientData } from '../../../../core/domain/dtos';
import { projectService } from '../../../project/services/project-service';

@Component({
  selector: 'app-clients-component',
  imports: [ListClientsComponent, UiSearchBarComponent, CreateButton],
  templateUrl: './clients-component.html',
  styleUrl: './clients-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ClientsComponent implements OnInit {
  protected clientService = inject(ClientService);
  private clientsMap = computed(() => this.clientService.state().clients);
  private dialog = inject(MatDialog);
  private projectService = inject(projectService)
  clientsList = computed(() => Array.from(this.clientsMap().values()));

  loading = computed(() => this.clientService.state().loading);

  ngOnInit(): void {
    this.clientService.getAllClients();
  }

  onSearch(query: string): void {
    this.clientService.searchClients(query);
  }

  openCreateClientModal(): void {
    const dialogRef = this.dialog.open(ClientFormComponent, {
      width: '400px',
      data: { client: null } 
    });

    dialogRef.afterClosed().subscribe(newClient => {
      if (newClient) {
        this.clientService.createClient(newClient);
      }
    });
  }

  openEditClientModal(clientToEdit: ClientData): void {
    const dialogRef = this.dialog.open(ClientFormComponent, {
      width: '400px',
      data: { client: clientToEdit } 
    });

    dialogRef.afterClosed().subscribe(updatedClient => {
      if (updatedClient) {
        this.clientService.updateClient(clientToEdit.id, updatedClient);
      }
    });
  }

  openCreateProjectModal(client: ClientData): void {
    const dialogRef = this.dialog.open(FormComponent, {
      width: '400px',
      data: { clientId: client.id }
    });

    dialogRef.afterClosed().subscribe(newProject => {
      if (newProject) {
        const projectPayload = {
          name: newProject.name,
          client_id: client.id
        };
        this.projectService.onCreate(projectPayload);
        this.clientService.getAllClients();
      }
    });
  }
}