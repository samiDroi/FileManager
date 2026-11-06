import { HttpInterceptorFn, HttpErrorResponse } from "@angular/common/http";
import { inject } from "@angular/core";
import { Router } from "@angular/router";
import { catchError, switchMap, throwError } from "rxjs";
import { AuthStateService } from "../../auth/auth-state.service";
import { AppRoutes } from "../routes/appRoutes";
import { ApiRoutes } from "../routes/apiUrls";
import { environment } from "../../../environments/environment";

export const authInterceptor: HttpInterceptorFn = (req, next) => {
  const authState = inject(AuthStateService);
  const router = inject(Router);
  const apiUrl = new URL(environment.apiUrl);
  const requestUrl = new URL(req.url, apiUrl);
  const apiPath = apiUrl.pathname.replace(/\/+$/, '');
  const isApiRequest = requestUrl.origin === apiUrl.origin &&
    (!apiPath || requestUrl.pathname === apiPath || requestUrl.pathname.startsWith(`${apiPath}/`));
  if (!isApiRequest) {
    return next(req);
  }

  const accessToken = authState.getAccessToken();
  const authReq = req.clone({
    withCredentials: true,
    setHeaders: accessToken ? { Authorization: `Bearer ${accessToken}` } : {}
  });
  const authEndpoints = [
    ApiRoutes.auth.login,
    ApiRoutes.auth.register,
    ApiRoutes.auth.logout,
    ApiRoutes.auth.refreshToken,
    ApiRoutes.auth.forgotPassword,
    ApiRoutes.auth.recoveryPassword
  ];
  const isAuthEndpoint = authEndpoints.includes(req.url);

  return next(authReq).pipe(
    catchError((error: HttpErrorResponse) => {
      if (error.status !== 401 || isAuthEndpoint || !authState.hasStoredUser()) {
        return throwError(() => error);
      }

      return authState.refreshAccessToken().pipe(
        switchMap(token => next(req.clone({
          withCredentials: true,
          setHeaders: { Authorization: `Bearer ${token}` }
        }))),
        catchError(refreshError => {
          authState.logout();
          router.navigate([AppRoutes.auth.login]);
          return throwError(() => refreshError);
        })
      );
    })
  );
};
