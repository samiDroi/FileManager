import { inject } from '@angular/core';
import { CanActivateChildFn, Router } from '@angular/router';
import { LocalKeys, LocalManagerService } from '../services/local-manager-service';
import { AppRoutes } from '../routes/appRoutes';
import { CheckPlatformUtils } from '../../utils';
import { AuthService } from '../../auth/auth-service';


export const authGuard: CanActivateChildFn = () => {
  // const localManager = inject(LocalManagerService);
  const router = inject(Router);
  const plattformUtils = inject(CheckPlatformUtils)
  const authService = inject(AuthService);

  if(plattformUtils.checkIfServer()){
    return true; // Si estamos en el servidor, permitimos la navegación sin verificar el token
  }

  if(authService.isLoggedIn()){
    return true; // Si el usuario está autenticado, permitimos la navegación
  }
  // if(!token){
  //   router.navigate([AppRoutes.auth.login]);
  // }

  // 3. Si definitivamente no hay sesión en el navegador, lo echamos.
  router.navigate([AppRoutes.auth.login]);
  return false;
};
