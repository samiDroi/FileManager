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
import { ClientService } from '../../clients/services';
import { Router } from '@angular/router';
import { AppRoutes } from '../../../core/routes/appRoutes';
import { PermissionService } from '../../../core/services/permission-service';

@Component({
  selector: 'app-projects-component',
  imports: [ListProjects, UiSearchBarComponent, CreateButton, EmptyStateComponent],
  templateUrl: './projects-component.html',
  styleUrl: './projects-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class projectsComponent implements OnInit {
  private projectService = inject(projectService);
  private clientService = inject(ClientService);
  private router = inject(Router);
  private permissionService = inject(PermissionService);
  protected modal = inject(MatDialog);

  protected readonly canCreateProjects = computed(() => this.permissionService.can('crear', 'proyectos'));
  protected readonly canEditProjects = computed(() => this.permissionService.can('editar', 'proyectos'));
  protected readonly canDeleteProjects = computed(() => this.permissionService.can('eliminar', 'proyectos'));

  projectsMap = computed(() => this.projectService.state().projects);
  projectsList = computed(() => Array.from(this.projectsMap().values()));
  loading = computed(() => this.projectService.state().loading);
  clientsLoading = computed(() => this.clientService.state().loading);
  hasClients = computed(() => this.clientService.state().clients.size > 0);

  ngOnInit(): void {
    this.projectService.getAll(); 
    this.clientService.getAllClients();
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

  onEmptyStateClick(): void {
    if (!this.hasClients()) {
      this.router.navigate(['/', AppRoutes.features.dashboard, AppRoutes.features.clients]);
      return;
    }

    this.openCreateProjectModal();
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

  onDeleteProject(project: ProjectData): void {
    void this.projectService.onDelete(project);
  }

  openFilesModal(project: ProjectData): void {
    this.modal.open(ProjectFilesModalComponent, {
      width: '850px',
      maxWidth: '95vw',
      data: { project: project }
    });
  }
}
