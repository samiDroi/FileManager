import { ChangeDetectionStrategy, Component, EventEmitter, Input, Output } from '@angular/core';
import { ProjectResponse } from '../../../../core';
import { UiTableComponent } from '../../../../shared/components/tables/ui-table-component/ui-table-component';

@Component({
  selector: 'app-list-projects',
  imports: [UiTableComponent],
  templateUrl: './list-projects.html',
  styleUrl: './list-projects.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class ListProjects {
  @Input({required: true}) projectsList: ProjectResponse[] = [];
  @Input() isLoading: boolean = false;

  @Output() editProject = new EventEmitter<ProjectResponse>();

  onEditItem(project: ProjectResponse): void {
    this.editProject.emit(project);
  }
}
