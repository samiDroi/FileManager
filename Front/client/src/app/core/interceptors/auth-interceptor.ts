import { HttpInterceptorFn, HttpErrorResponse } from "@angular/common/http";
import { inject } from "@angular/core";
import { Router } from "@angular/router";
import { catchError, throwError } from "rxjs";
import { AuthStateService } from "../../auth/auth-state.service";
import { AppRoutes } from "../routes/appRoutes";

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authState = inject(AuthStateService);
  const router = inject(Router);

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
