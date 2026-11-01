import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { toSignal } from '@angular/core/rxjs-interop';
import { NavigationEnd, Router, RouterOutlet } from '@angular/router';
import { filter } from 'rxjs';
import { AuthService } from '../../auth/auth-service';
import { SidebarComponent } from './sidebar-component/sidebar-component';
import { AppRoutes } from '../../core/routes/appRoutes';
import { TopbarTitleService } from '../../core/services';
import { ClientService } from '../../features/clients/services';
import { projectService } from '../../features/project/services';
import { MatDialog } from '@angular/material/dialog';
import { ClientFormComponent } from '../../features/clients/components/clients-component/client-form-component/client-form-component';

@Component({
  selector: 'app-main-layout',
  imports: [SidebarComponent, RouterOutlet],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MainLayout implements OnInit {
  private authService = inject(AuthService);
  private clientService = inject(ClientService)
  private projectService = inject(projectService);
  private topbarTitleService = inject(TopbarTitleService);
  private router = inject(Router);
  private dialog = inject(MatDialog);
  protected currentUser = this.authService.currentUser;
  private navigationEnd = toSignal(
    this.router.events.pipe(
      filter((event): event is NavigationEnd => event instanceof NavigationEnd)
    ),
    { initialValue: null }
  );
  private currentUrl = computed(() =>
    this.navigationEnd()?.urlAfterRedirects ?? this.router.url
  );

  clients = computed(() => Array.from(this.clientService.state().clients.values()));

  protected pageTitle = computed(() => {
    const projectContext = this.topbarTitleService.projectContext();

    if (projectContext) {
      const projectTitle = [projectContext.clientName, projectContext.projectName]
        .filter(Boolean)
        .join('/');

      if (projectTitle) {
        return projectTitle;
      }
    }

    return this.getRouteTitle();
  });

  private getRouteTitle(): string {
    const segments = this.currentUrl().split('?')[0].split('/').filter(Boolean);
    const feature = segments[1];

    if (!feature || feature === AppRoutes.features.dashboard) {
      return 'Bienvenido de nuevo';
    }

    if (feature === 'client') {
      const clientId = Number(segments[2]);
      const client =
        this.clientService.state().clients.get(clientId) ??
        this.clientService.stateWithProjects().clientsWithProjects.find(c => c.id === clientId);
      const projectIndex = segments.indexOf('project');

      if (projectIndex !== -1) {
        const projectId = Number(segments[projectIndex + 1]);
        const project = this.projectService.state().projects.get(projectId);
        const projectTitle = [client?.name, project?.name].filter(Boolean).join('/');

        if (projectTitle) {
          return projectTitle;
        }
      }

      return client?.name || 'Cliente';
    }

    const sectionTitles: Record<string, string> = {
      clients: 'Clientes',
      projects: 'Proyectos',
      files: 'Archivos',
      users: 'Usuarios',
    };

    return sectionTitles[feature] ?? feature;
  }

  ngOnInit(): void {
    this.clientService.getAllClients();
  }

  openCreateClientModal(): void {
    const dialogRef = this.dialog.open(ClientFormComponent, {
      width: '400px',
      data: { client: null }
    });

    dialogRef.afterClosed().subscribe(newClient => {
      if (newClient) {
        this.clientService.createClient(newClient);
      }
    });
  }

  handleLogout(): void {
    this.authService.logout();
    this.router.navigate([AppRoutes.auth.login]);
  }

}
