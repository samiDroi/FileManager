import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { map } from 'rxjs';
import { Location } from '@angular/common';
import { ProjectData } from '../../../core';
import { UiTableComponent } from '../../../shared/components/tables/ui-table-component/ui-table-component';
import { EditButton } from '../../../shared/components/buttons/edit-button/edit-button';
import { CreateButton } from '../../../shared/components/buttons/create-button/create-button';
import { projectService } from '../services';
import { ClientService } from '../../../features/clients/services';
import { MatDialog } from '@angular/material/dialog';
import { FormComponent } from '../projects-component/form-component/form-component';
import { ProjectFilesModalComponent } from '../../files/components/project-files-modal';
import { UiSearchBarComponent } from '../../../shared/components/search-bar/ui-search-bar-component';

@Component({
  selector: 'app-dynamic-list',
  imports: [UiTableComponent, EditButton, CreateButton,UiSearchBarComponent],
  templateUrl: './dynamic-list-component.html',
  styleUrl: './dynamic-list-component.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class DynamicListComponent {
  private route = inject(ActivatedRoute);
  private location = inject(Location);
  private projectService = inject(projectService);
  private clientService = inject(ClientService);
  protected modal = inject(MatDialog);

  clientId = toSignal(
    this.route.params.pipe(map(p => Number(p['id']))),
    { initialValue: 0 }
  );

  clientName = computed(() => {
    const id = this.clientId();
    if (!id) return '';
    const clients = this.clientService.stateWithProjects().clientsWithProjects;
    const client = clients.find(c => c.id === id);
    return client?.name ?? '';
  });

  allProjects = computed(() => Array.from(this.projectService.state().projects.values()));
  loading = computed(() => this.projectService.state().loading);

  projects = computed(() => {
    const id = this.clientId();
    if (!id) return [];
    //FILTER PROJECTS BY CLIENT ID, CAUSE THE PROJECTDATA ALREADY HAS THE CLIENT ID
    return this.allProjects().filter(p => p.client_id === id);
  });

  constructor() {
    this.projectService.getAll();
    this.clientService.getClientsWithProjects();
  }

  goBack(): void {
    this.location.back();
  }

  openCreateProjectModal(): void {
    const dialogRef = this.modal.open(FormComponent, {
      width: '500px',
      data: { project: null, preselectedClientId: this.clientId() }
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

  onEditProject(project: ProjectData): void {
    const dialogRef = this.modal.open(FormComponent, {
      width: '500px',
      data: { project }
    });

    dialogRef.afterClosed().subscribe(updated => {
      if (updated) {
        this.projectService.onUpdate(updated.id, updated);
        this.projectService.getAll();
      }
    });
  }

  onGoToFiles(project: ProjectData): void {
    this.modal.open(ProjectFilesModalComponent, {
      width: '850px',
      maxWidth: '95vw',
      data: { project }
    });
  }
  onSearch(query: string): void {
    this.projectService.searchProjects(query);
  }
}
