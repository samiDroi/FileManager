import { ChangeDetectionStrategy, Component, input, output } from '@angular/core';
import { UserData } from '../../../../../core/domain/dtos/user';
import { UiTableComponent } from '../../../../../shared/components/tables/ui-table-component/ui-table-component';
import { EditButton } from '../../../../../shared/components/buttons/edit-button/edit-button';
import { DeleteButton } from '../../../../../shared/components/buttons/delete-button/delete-button';

@Component({
  selector: 'app-list-users',
  imports: [UiTableComponent,EditButton,DeleteButton],
  templateUrl: './list-users-component.html',
  styleUrl: './list-users-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListUsersComponent {
  usersList = input<UserData[]>([]);
  canEdit = input<boolean>(false);
  canDelete = input<boolean>(false);
  deleteUser = output<UserData>();
  editUser = output<UserData>();
  onEdit(user: UserData) {
    this.editUser.emit(user);
  }
  onDelete(user: UserData) {
    this.deleteUser.emit(user);
  }
}
