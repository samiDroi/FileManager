import { HttpClient } from '@angular/common/http';
import { inject, Injectable } from '@angular/core';
import { finalize, Observable } from 'rxjs';
import { AuthStateService } from './auth-state.service';
import { ApiRoutes } from '../core/routes/apiUrls';
import { Router } from '@angular/router';
import { AppRoutes } from '../core/routes/appRoutes';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private http = inject(HttpClient);
  private authState = inject(AuthStateService);
  private router = inject(Router)
  readonly currentUser = this.authState.currentUser;
  readonly isLoggedIn = this.authState.isLoggedIn;

  login(credentials: any): Observable<any> {
    return this.authState.login(credentials);
  }

  register(data: any): Observable<any> {
    return this.authState.register(data);
  }

  logout(): void {
    this.http.post(ApiRoutes.auth.logout, {}, { withCredentials: true })
    .pipe(
      finalize(()=>{
        this.authState.logout()
        this.router.navigate([AppRoutes.auth.login])
      })
    ).subscribe();
  }

}
