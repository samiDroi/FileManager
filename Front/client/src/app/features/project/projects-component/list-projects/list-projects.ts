import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output, input } from '@angular/core';
import { ProjectData } from '../../../../core';
import { UiTableComponent } from '../../../../shared/components/tables/ui-table-component/ui-table-component';
import { EditButton } from '../../../../shared/components/buttons/edit-button/edit-button';
import { DeleteButton } from '../../../../shared/components/buttons/delete-button/delete-button';

@Component({
  selector: 'app-list-projects',
  imports: [UiTableComponent, EditButton, DeleteButton],
  templateUrl: './list-projects.html',
  styleUrl: './list-projects.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListProjects {
  projectsList = input.required<ProjectData[]>();
  isLoading = input<boolean>(false);
  canEdit = input<boolean>(false);
  canDelete = input<boolean>(false);

  @Output() editProject = new EventEmitter<ProjectData>();
  @Output() deleteProject = new EventEmitter<ProjectData>();
  @Output() goToFiles = new EventEmitter<ProjectData>();

  onEditItem(project: ProjectData): void {
    this.editProject.emit(project);
  }

  onGoToFiles(project: ProjectData): void {
    this.goToFiles.emit(project);
  }

  onDeleteItem(project: ProjectData): void {
    this.deleteProject.emit(project);
  }
}