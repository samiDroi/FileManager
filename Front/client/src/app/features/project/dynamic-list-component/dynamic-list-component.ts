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
import { EmptyStateComponent } from '../../../shared/components';

@Component({
  selector: 'app-dynamic-list',
  imports: [UiTableComponent, EditButton, CreateButton, UiSearchBarComponent, EmptyStateComponent],
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
    this.route.paramMap.pipe(
      map(params => {
        const id = Number(params.get('id'));
        return Number.isInteger(id) && id > 0 ? id : null;
      })
    ),
    { initialValue: null }
  );

  clientName = computed(() => {
    const id = this.clientId();
    if (id === null) return '';
    const clients = this.clientService.stateWithProjects().clientsWithProjects;
    const client = clients.find(c => c.id === id);
    return client?.name ?? '';
  });

  allProjects = computed(() => Array.from(this.projectService.state().projects.values()));
  loading = computed(() => this.projectService.state().loading);
  searchQuery = signal('');

  protected clientProjects = computed(() => {
    const id = this.clientId();
    if (id === null) return [];
    return this.allProjects().filter(project => Number(project.client_id) === id);
  });

  projects = computed(() => {
    const query = this.searchQuery().trim().toLocaleLowerCase();
    if (!query) return this.clientProjects();
    return this.clientProjects().filter(project =>
      `${project.name} ${project.client_name}`.toLocaleLowerCase().includes(query)
    );
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
    this.searchQuery.set(query);
  }
}
