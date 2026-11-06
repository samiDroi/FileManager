import { TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { firstValueFrom, Observable, of } from 'rxjs';
import { PermissionService } from '../services/permission-service';
import { AppRoutes } from '../routes/appRoutes';
import { adminGuard } from './admin-guard';

describe('adminGuard', () => {
  let router: { createUrlTree: ReturnType<typeof vi.fn> };
  let permissionService: {
    isAdmin: ReturnType<typeof vi.fn>;
    loaded: ReturnType<typeof vi.fn>;
    load: ReturnType<typeof vi.fn>;
  };

  const runGuard = () =>
    TestBed.runInInjectionContext(() => adminGuard({} as never, {} as never));

  beforeEach(() => {
    router = {
      createUrlTree: vi.fn(() => ({ redirected: true })),
    };
    permissionService = {
      isAdmin: vi.fn(() => true),
      loaded: vi.fn(() => true),
      load: vi.fn(() => of(true)),
    };

    TestBed.configureTestingModule({
      providers: [
        { provide: Router, useValue: router },
        { provide: PermissionService, useValue: permissionService },
      ],
    });
  });

  it('allows administrators to access admin routes', () => {
    expect(runGuard()).toBe(true);
  });

  it('redirects non-administrators to the main dashboard', () => {
    permissionService.isAdmin.mockReturnValue(false);

    expect(runGuard()).toEqual({ redirected: true });
    expect(router.createUrlTree).toHaveBeenCalledWith([AppRoutes.features.dashboard]);
  });

  it('loads permissions before deciding when they are not available yet', async () => {
    permissionService.isAdmin.mockReturnValue(false);
    permissionService.loaded.mockReturnValue(false);
    permissionService.load.mockImplementation(() => {
      permissionService.isAdmin.mockReturnValue(true);
      return of(true);
    });

    const result = runGuard() as Observable<boolean>;

    await expect(firstValueFrom(result)).resolves.toBe(true);
    expect(permissionService.load).toHaveBeenCalled();
  });
});
