import { ChangeDetectionStrategy, Component, inject } from '@angular/core';
import { AuthService } from '../../auth/auth-service';
import { SidebarComponent } from './sidebar-component/sidebar-component';
import { Router, RouterOutlet } from '@angular/router';
import { AppRoutes } from '../../core/routes/appRoutes';

@Component({
  selector: 'app-main-layout',
  imports: [SidebarComponent, RouterOutlet],
  templateUrl: './main-layout.html',
  styleUrl: './main-layout.scss',
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class MainLayout {
  private authService = inject(AuthService);
  private router = inject(Router);
  protected currentUser = this.authService.currentUser;

  handleLogout(): void {
    this.authService.logout();
    this.router.navigate([AppRoutes.auth.login]);
  }
}
