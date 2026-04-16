import { ChangeDetectionStrategy, Component, Input, Output, EventEmitter, output } from '@angular/core';
import { Client, Proyect } from '../../../../../core';
import { UiTableComponent } from '../../../../../shared/components/tables/ui-table-component/ui-table-component';

@Component({
  selector: 'app-list-clients',
  imports: [UiTableComponent],
  templateUrl: './list-clients-component.html',
  styleUrl: './list-clients-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListClientsComponent {
  @Input({required: true}) clients: Client[] = [];
  @Input() isLoading: boolean = false;

  @Output() createClient = new EventEmitter<void>();
  @Output() editClient = new EventEmitter<Client>();
  @Output() createProyect = new EventEmitter<Proyect>();

  onCreateItem(client: Client): void {
    this.createClient.emit();
  }

  onEditItem(client: Client): void {
    this.editClient.emit(client);
  }

  
}
