import { ChangeDetectionStrategy, Component, computed, inject, OnInit } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { AuthService } from '../../auth/auth-service';
import { AppRoutes } from '../../core/routes/appRoutes';
import { PermissionService } from '../../core/services';
import { ClientService } from '../../features/clients/services';

@Component({
  selector: 'app-admin-layout',
  imports: [RouterLink, RouterLinkActive, RouterOutlet],
  templateUrl: './admin-layout.html',
  styleUrl: './admin-layout.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdminLayout implements OnInit {
  protected readonly routes = AppRoutes.admin;
  private readonly authService = inject(AuthService);
  private readonly permissionService = inject(PermissionService);
  private readonly clientService = inject(ClientService);
  protected readonly canManageRbac = computed(() => this.permissionService.can('gestionar', 'roles'));
  protected readonly canReadUsers = computed(() => this.permissionService.can('leer', 'usuarios'));
  protected readonly currentUser = this.authService.currentUser;

  ngOnInit(): void {
    this.permissionService.load().subscribe();
    this.clientService.getAllClients();
  }

  protected handleLogout(): void {
    this.authService.logout();
  }
}
