import { ChangeDetectionStrategy, Component, Output, EventEmitter, input } from '@angular/core';
import { ClientData } from '../../../../../core/domain/dtos';
import { UiTableComponent } from '../../../../../shared/components/tables/ui-table-component/ui-table-component';
import { CreateButton } from '../../../../../shared/components/buttons/create-button/create-button';
import { EditButton } from '../../../../../shared/components/buttons/edit-button/edit-button';
import { DeleteButton } from '../../../../../shared/components/buttons/delete-button/delete-button';

@Component({
  selector: 'app-list-clients',
  imports: [UiTableComponent, CreateButton, EditButton, DeleteButton],
  templateUrl: './list-clients-component.html',
  styleUrl: './list-clients-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListClientsComponent {
  clients = input.required<ClientData[]>();
  isLoading = input<boolean>(false);
  canCreateProject = input<boolean>(false);
  canEdit = input<boolean>(false);
  canDelete = input<boolean>(false);

  @Output() createClient = new EventEmitter<void>();
  @Output() editClient = new EventEmitter<ClientData>();
  @Output() deleteClient = new EventEmitter<ClientData>();
  @Output() createproject = new EventEmitter<ClientData>();

  onCreateProject(client: ClientData): void {
    this.createproject.emit(client);
  }

  onEditItem(client: ClientData): void {
    this.editClient.emit(client);
  }

  onDeleteItem(client: ClientData): void {
    this.deleteClient.emit(client);
  }
}