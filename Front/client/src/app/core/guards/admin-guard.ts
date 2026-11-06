import { inject } from '@angular/core';
import { CanActivateFn, Router } from '@angular/router';
import { map } from 'rxjs';
import { AppRoutes } from '../routes/appRoutes';
import { PermissionService } from '../services/permission-service';

export const adminGuard: CanActivateFn = () => {
  const router = inject(Router);
  const permissionService = inject(PermissionService);

  const redirectToMain = () => router.createUrlTree([AppRoutes.features.dashboard]);

  if (permissionService.isAdmin()) {
    return true;
  }

  if (permissionService.loaded()) {
    return redirectToMain();
  }

  return permissionService
    .load()
    .pipe(map(() => (permissionService.isAdmin() ? true : redirectToMain())));
};
