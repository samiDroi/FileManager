import { HttpInterceptorFn, HttpErrorResponse } from "@angular/common/http";
import { inject, PLATFORM_ID } from "@angular/core";
import { Router } from "@angular/router";
import { catchError, EMPTY, throwError } from "rxjs";
import { AuthStateService } from "../../auth/auth-state.service";
import { AppRoutes } from "../routes/appRoutes";
import { isPlatformBrowser } from "@angular/common";

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  
  const authState = inject(AuthStateService);
  const router = inject(Router);
  // 1. Si NO estamos en el navegador (estamos en Node.js), cancelamos la petición
  // if (!isPlatformBrowser(PLATFORM_ID)) {
  //   console.log('Frenando petición en el servidor SSR:', req.url);
  //   return EMPTY; // Corta la ejecución aquí, no hace la petición HTTP
  // }
  const authReq = req.clone({
    withCredentials: true
  });

  return next(authReq).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status === 401) {
        if (authState.hasStoredUser()) {
          authState.logout();
          router.navigate([AppRoutes.auth.login]);
        }
      }
      return throwError(() => error);
    })
  );
};
