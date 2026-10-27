import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { AuthService } from '../../auth/auth-service';
import { SidebarComponent } from './sidebar-component/sidebar-component';
import { Router, RouterOutlet } from '@angular/router';
import { AppRoutes } from '../../core/routes/appRoutes';
import { ClientService } from '../../features/clients/services';
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
  private router = inject(Router);
  private dialog = inject(MatDialog);
  protected currentUser = this.authService.currentUser;
  
  clients = computed(() => Array.from(this.clientService.state().clients.values()));

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
