import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { of } from 'rxjs';
import { AuthService } from '../../auth-service';
import { AppRoutes } from '../../../core/routes/appRoutes';

import { LoginComponent } from './login-component';

describe('LoginComponent', () => {
  let component: LoginComponent;
  let fixture: ComponentFixture<LoginComponent>;
  let router: { navigate: ReturnType<typeof vi.fn> };
  let authService: { login: ReturnType<typeof vi.fn> };

  beforeEach(async () => {
    router = { navigate: vi.fn() };
    authService = {
      login: vi.fn(() => of({ is_admin: false })),
    };

    await TestBed.configureTestingModule({
      imports: [LoginComponent],
      providers: [
        { provide: Router, useValue: router },
        { provide: AuthService, useValue: authService },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(LoginComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('redirects administrators to the admin dashboard', async () => {
    authService.login.mockReturnValue(of({ is_admin: true }));
    component.loginForm.setValue({ username: 'admin', password: 'secret' });

    await component.onSubmit();

    expect(router.navigate).toHaveBeenCalledWith([AppRoutes.admin.dashboard]);
  });

  it('redirects regular users to the main dashboard', async () => {
    authService.login.mockReturnValue(of({ is_admin: false }));
    component.loginForm.setValue({ username: 'user', password: 'secret' });

    await component.onSubmit();

    expect(router.navigate).toHaveBeenCalledWith([AppRoutes.features.dashboard]);
  });
});
