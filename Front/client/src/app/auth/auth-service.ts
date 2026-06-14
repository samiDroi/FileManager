import { HttpClient } from '@angular/common/http';
import { inject, Injectable, OnInit } from '@angular/core';
import { finalize, Observable } from 'rxjs';
import { AuthStateService } from './auth-state.service';
import { ApiRoutes } from '../core/routes/apiUrls';
import { ActivatedRoute, Router } from '@angular/router';
import { AppRoutes } from '../core/routes/appRoutes';
import Swal from 'sweetalert2';

@Injectable({
  providedIn: 'root',
})
export class AuthService  {
  private http = inject(HttpClient);
  
  private authState = inject(AuthStateService);
  private router = inject(Router)
  
  readonly currentUser = this.authState.currentUser;
  

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
  //here i can't use one source of truth because the user is not logged in yet, so i need to use the http client directly
  sendResetPasswordEmail(email: string): Observable<any> {
    return this.http.post(ApiRoutes.auth.forgotPassword, { email });
  }

  recoveryPassword(token: string, newPassword: string): Observable<any> {
    return this.http.post(ApiRoutes.auth.recoveryPassword, { token, new_password: newPassword });
  }
}
