import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { AuthService } from '../../auth/auth-service';
import { AppRoutes } from '../../core/routes/appRoutes';

@Component({
  selector: 'app-admin-layout',
  imports: [RouterLink, RouterLinkActive, RouterOutlet],
  templateUrl: './admin-layout.html',
  styleUrl: './admin-layout.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class AdminLayout {
  protected readonly routes = AppRoutes.admin;
  private readonly authService = inject(AuthService);
  protected readonly currentUser = this.authService.currentUser;

  protected handleLogout(): void {
    this.authService.logout();
  }
}
