import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { ListUsersComponent } from './list-users-component/list-users-component';
import { UserService } from '../../services/user-service';
import { UiSearchBarComponent } from '../../../../shared/components/search-bar/ui-search-bar-component';
import { UserData } from '../../../../core/domain/dtos/user';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
import { MatDialog } from '@angular/material/dialog';
import { RegisterComponent } from '../../../../auth';
import { EmptyStateComponent } from '../../../../shared/components';

@Component({
  selector: 'app-user-component',
  imports: [ListUsersComponent, UiSearchBarComponent, CreateButton, EmptyStateComponent],
  templateUrl: './user-component.html',
  styleUrl: './user-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class UserComponent implements OnInit {
  private userService =  inject(UserService);
  private modal = inject(MatDialog);
  users = computed(() => Array.from(this.userService.state().users.values()));

  ngOnInit() {
    
    this.userService.getAllUsers();
    console.log(this.users());
  }
  
  onSearch(query: string) {
      this.userService.searchUsers(query);
  }

  onEdit(user: UserData) {
    this.userService.editUser(user);
  }
  onDelete(user: UserData) {
    this.userService.deleteUser(user.id);
  }
  openCreateUserModal() {
    this.modal.open(RegisterComponent, {
      width: '70%',
      panelClass: 'custom-modal',
    });
  }
}

