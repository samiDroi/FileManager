import { ChangeDetectionStrategy, computed, inject, OnInit,Component } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { FormComponent } from './form-component/form-component';
import { ListProjects } from './list-projects/list-projects';
import { UiSearchBarComponent } from '../../../shared/components/search-bar/ui-search-bar-component';
import { CreateButton } from '../../../shared/components/buttons/create-button/create-button';
import { projectService } from '../services';
import { ProjectData } from '../../../core';
import { ProjectFilesModalComponent } from '../../files/components/project-files-modal';
import { EmptyStateComponent } from '../../../shared/components';

@Component({
  selector: 'app-projects-component',
  imports: [ListProjects, UiSearchBarComponent, CreateButton, EmptyStateComponent],
  templateUrl: './projects-component.html',
  styleUrl: './projects-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class projectsComponent implements OnInit {
  private projectService = inject(projectService);
  protected modal = inject(MatDialog);

  projectsMap = computed(() => this.projectService.state().projects);
  projectsList = computed(() => Array.from(this.projectsMap().values()));
  loading = computed(() => this.projectService.state().loading);

  ngOnInit(): void {
    this.projectService.getAll(); 
  }

  onSearch(query: string): void {
    this.projectService.searchProjects(query);
  }

  openCreateProjectModal(): void {
    const dialogRef = this.modal.open(FormComponent, {
      width: '500px',
      data: { project: null }
    });

    dialogRef.afterClosed().subscribe(newProject => {
      if (newProject) {
        this.projectService.onCreate({
          name: newProject.name,
          client_id: newProject.client_id
        });
      }
    });
  }

  openEditProjectModal(project: ProjectData): void {
    const dialogRef = this.modal.open(FormComponent, {
      width: '500px',
      data: { project: project }
    });

    dialogRef.afterClosed().subscribe(updatedProject => {
      if (updatedProject) {
        console.log('Actualizar proyecto:', updatedProject);
        this.projectService.onUpdate(updatedProject.id, updatedProject);
      }
    });
  }

  openFilesModal(project: ProjectData): void {
    this.modal.open(ProjectFilesModalComponent, {
      width: '850px',
      maxWidth: '95vw',
      data: { project: project }
    });
  }
}
