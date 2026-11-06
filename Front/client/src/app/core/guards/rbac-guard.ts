import { inject } from '@angular/core';
import { CanActivateChildFn, Router } from '@angular/router';
import { map } from 'rxjs';
import { PermissionService } from '../services/permission-service';
import { AppRoutes } from '../routes/appRoutes';

export const rbacGuard: CanActivateChildFn = () => {
  const router = inject(Router);
  const permissionService = inject(PermissionService);

  return permissionService
    .load()
    .pipe(
      map(canManage =>
        canManage ? true : router.createUrlTree([AppRoutes.admin.dashboard])
      )
    );
};
