import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { GroupData, GroupRequest } from '../../../../core/domain/dtos/rbac';
import { EmptyStateComponent } from '../../../../shared/components';
import { CreateButton } from '../../../../shared/components/buttons/create-button/create-button';
import { GroupFormDialog } from '../group-form-dialog/group-form-dialog';
import { GroupDetailDialog } from '../group-detail-dialog/group-detail-dialog';
import { GroupService } from '../../services';
import { PermissionService } from '../../../../core/services/permission-service';

@Component({
  selector: 'app-groups-component',
  imports: [CreateButton, EmptyStateComponent],
  templateUrl: './groups-component.html',
  styleUrl: './groups-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class GroupsComponent implements OnInit {
  private readonly groupService = inject(GroupService);
  private readonly permissionService = inject(PermissionService);
  private readonly dialog = inject(MatDialog);

  protected readonly groups = this.groupService.groups;
  protected readonly loading = this.groupService.loading;
  protected readonly canManageGroups = computed(() => this.permissionService.can('gestionar', 'roles'));

  ngOnInit(): void {
    this.groupService.getAllGroups();
  }

  protected openCreateGroup(): void {
    const dialogRef = this.dialog.open<GroupFormDialog, { group: GroupData | null }, GroupRequest>(
      GroupFormDialog,
      {
        width: '460px',
        maxWidth: '95vw',
        data: { group: null },
      }
    );

    dialogRef.afterClosed().subscribe(request => {
      if (request) {
        this.groupService.createGroup(request);
      }
    });
  }

  protected openEditGroup(group: GroupData): void {
    const dialogRef = this.dialog.open<GroupFormDialog, { group: GroupData | null }, GroupRequest>(
      GroupFormDialog,
      {
        width: '460px',
        maxWidth: '95vw',
        data: { group },
      }
    );

    dialogRef.afterClosed().subscribe(request => {
      if (request) {
        this.groupService.updateGroup(group.id, request);
      }
    });
  }

  protected openGroupDetail(group: GroupData): void {
    this.dialog.open(GroupDetailDialog, {
      width: '940px',
      maxWidth: '95vw',
      maxHeight: '92vh',
      data: { group },
    });
  }

  protected deleteGroup(group: GroupData, event: Event): void {
    event.stopPropagation();
    void this.groupService.deleteGroup(group);
  }
}
