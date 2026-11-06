import { ChangeDetectionStrategy, Component, computed, inject, OnInit, signal } from '@angular/core';
import { ListUsersComponent } from './list-users-component/list-users-component';
import { UserService } from '../../services/user-service';
import { UiSearchBarComponent } from '../../../../shared/components/search-bar/ui-search-bar-component';
import { UserData } from '../../../../core/domain/dtos/user';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
import { MatDialog } from '@angular/material/dialog';
import { EmptyStateComponent } from '../../../../shared/components';
import { CreateUserDialog } from '../create-user-dialog/create-user-dialog';
import { PermissionService } from '../../../../core/services/permission-service';
import { EditUserDialog } from '../edit-user-dialog/edit-user-dialog';

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
  private permissionService = inject(PermissionService);
  canCreateUsers = computed(() => this.permissionService.can('crear', 'usuarios'));
  canEditUsers = computed(() => this.permissionService.can('editar', 'usuarios'));
  canDeleteUsers = computed(() => this.permissionService.can('eliminar', 'usuarios'));
  users = computed(() => Array.from(this.userService.state().users.values()));
  successMessage = signal('');

  ngOnInit() {
    this.userService.getAllUsers();
  }
  
  onSearch(query: string) {
    this.userService.searchUsers(query);
  }

  onEdit(user: UserData) {
    this.successMessage.set('');
    this.modal.open(EditUserDialog, {
      width: '70%',
      panelClass: 'custom-modal',
      data: user,
    }).afterClosed().subscribe(updatedUser => {
      if (updatedUser) {
        this.userService.upsertUser(updatedUser);
        this.successMessage.set(`El usuario ${updatedUser.username} se actualizó.`);
      }
    });
  }
  onDelete(user: UserData) {
    this.userService.deleteUser(user.id);
  }
  openCreateUserModal() {
    this.successMessage.set('');
    this.modal.open(CreateUserDialog, {
      width: '70%',
      panelClass: 'custom-modal',
    }).afterClosed().subscribe(createdUser => {
      if (createdUser) {
        this.userService.upsertUser(createdUser);
        this.successMessage.set(`El usuario ${createdUser.username} se dio de alta.`);
      }
    });
  }
}
