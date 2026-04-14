import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { ClientService } from '../../services/client-service';
import { ListClientsComponent } from './list-clients-component/list-clients-component';
import { UiSearchBarComponent } from '../../../../shared/components/search-bar/ui-search-bar-component';
import { MatDialog } from '@angular/material/dialog';
import { ClientFormComponent } from './client-form-component/client-form-component';
import { Client } from '../../../../core';

@Component({
  selector: 'app-clients-component',
  imports: [ListClientsComponent, UiSearchBarComponent],
  templateUrl: './clients-component.html',
  styleUrl: './clients-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ClientsComponent implements OnInit {
  protected clientService = inject(ClientService);
  private clientsMap = computed(() => this.clientService.state().clients);
  private dialog = inject(MatDialog);
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
        console.log('Crear cliente:', newClient);
        this.clientService.createClient(newClient);
      }
    });
  }

  openEditClientModal(clientToEdit: Client): void {
    const dialogRef = this.dialog.open(ClientFormComponent, {
      width: '400px',
      data: { client: clientToEdit } 
    });

    dialogRef.afterClosed().subscribe(updatedClient => {
      if (updatedClient) {
        console.log('Actualizar cliente:', updatedClient);
        this.clientService.updateClient(clientToEdit.id, updatedClient);
      }
    });
  }
}
