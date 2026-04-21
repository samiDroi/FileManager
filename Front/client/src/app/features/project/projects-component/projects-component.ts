import { ChangeDetectionStrategy, computed, inject, OnInit,Component } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { FormComponent } from './form-component/form-component';
import { ListProjects } from './list-projects/list-projects';
import { UiSearchBarComponent } from '../../../shared/components/search-bar/ui-search-bar-component';
import { projectService } from '../services';
import { ProjectResponse } from '../../../core';

@Component({
  selector: 'app-projects-component',
  imports: [ListProjects, UiSearchBarComponent],
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
        this.projectService.getAll();
      }
    });
  }

  openEditProjectModal(project: ProjectResponse): void {
    const dialogRef = this.modal.open(FormComponent, {
      width: '500px',
      data: { project: project }
    });

    dialogRef.afterClosed().subscribe(updatedProject => {
      if (updatedProject) {
        console.log('Actualizar proyecto:', updatedProject);
        this.projectService.onUpdate(updatedProject.id, updatedProject);
        this.projectService.getAll();
      }
    });
  }
}
