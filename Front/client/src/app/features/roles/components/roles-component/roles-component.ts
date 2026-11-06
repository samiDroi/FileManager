import { ChangeDetectionStrategy, Component, inject, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { RoleData, RoleRequest } from '../../../../core/domain/dtos/rbac';
import { EmptyStateComponent } from '../../../../shared/components';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
import { RoleFormDialog, RoleFormDialogData } from '../role-form-dialog/role-form-dialog';
import { RoleService } from '../../services';

@Component({
  selector: 'app-roles-component',
  imports: [CreateButton, EmptyStateComponent],
  templateUrl: './roles-component.html',
  styleUrl: './roles-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class RolesComponent implements OnInit {
  private readonly roleService = inject(RoleService);
  private readonly dialog = inject(MatDialog);

  protected readonly roles = this.roleService.roles;
  protected readonly permissions = this.roleService.permissions;
  protected readonly loading = this.roleService.loading;

  ngOnInit(): void {
    this.roleService.getAllRoles();
    this.roleService.getAllPermissions();
  }

  protected openCreateRole(): void {
    const data: RoleFormDialogData = {
      role: null,
      permissions: this.permissions(),
    };
    const dialogRef = this.dialog.open<RoleFormDialog, RoleFormDialogData, RoleRequest>(
      RoleFormDialog,
      {
        width: '620px',
        maxWidth: '95vw',
        maxHeight: '92vh',
        data,
      }
    );

    dialogRef.afterClosed().subscribe(request => {
      if (request) {
        this.roleService.createRole(request);
      }
    });
  }

  protected openEditRole(role: RoleData): void {
    const data: RoleFormDialogData = {
      role,
      permissions: this.permissions(),
    };
    const dialogRef = this.dialog.open<RoleFormDialog, RoleFormDialogData, RoleRequest>(
      RoleFormDialog,
      {
        width: '620px',
        maxWidth: '95vw',
        maxHeight: '92vh',
        data,
      }
    );

    dialogRef.afterClosed().subscribe(request => {
      if (request) {
        this.roleService.updateRole(role.id, request);
      }
    });
  }

  protected deleteRole(role: RoleData): void {
    void this.roleService.deleteRole(role);
  }

  protected isSystemRole(role: RoleData): boolean {
    const normalizedName = role.name.trim().toLowerCase();
    return normalizedName === 'administrador' || normalizedName === 'admin';
  }
}
