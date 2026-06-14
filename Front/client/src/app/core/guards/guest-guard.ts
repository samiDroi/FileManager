import { inject } from '@angular/core';
import { CanActivateChildFn, Router } from '@angular/router';
import { AppRoutes } from '../routes/appRoutes';
import { CheckPlatformUtils } from '../../utils';
import { AuthStateService } from '../../auth/auth-state.service';

export const guestGuard: CanActivateChildFn = () => {
  //guard for routes that should only be accessible to guests (not logged in users)
  const router = inject(Router);
  const platformUtils = inject(CheckPlatformUtils);
  const authState = inject(AuthStateService);

  if (platformUtils.checkIfServer()) {
    return true;
  }

  if (!authState.currentUser()) {
    return true;
  }

  router.navigate([AppRoutes.features.dashboard]);
  return false;
};
