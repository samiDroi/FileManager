import { inject } from '@angular/core';
import { CanActivateChildFn, Router } from '@angular/router';
import { AppRoutes } from '../routes/appRoutes';
import { CheckPlatformUtils } from '../../utils';
import { AuthStateService } from '../../auth/auth-state.service';
import { PermissionService } from '../services/permission-service';

export const guestGuard: CanActivateChildFn = () => {
  //guard for routes that should only be accessible to guests (not logged in users)
  const router = inject(Router);
  const platformUtils = inject(CheckPlatformUtils);
  const authState = inject(AuthStateService);
  const permissionService = inject(PermissionService);

  if (platformUtils.checkIfServer()) {
    return true;
  }

  if (!authState.currentUser()) {
    return true;
  }

  router.navigate([
    permissionService.isAdmin()
      ? AppRoutes.admin.dashboard
      : AppRoutes.features.dashboard,
  ]);
  return false;
};
