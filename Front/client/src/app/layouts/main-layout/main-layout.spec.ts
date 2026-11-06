import { ComponentFixture, TestBed } from '@angular/core/testing';
import { signal } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { provideRouter } from '@angular/router';
import { of } from 'rxjs';
import { AuthService } from '../../auth/auth-service';
import { TopbarTitleService } from '../../core/services';
import { PermissionService } from '../../core/services/permission-service';
import { ClientService } from '../../features/clients/services';
import { projectService } from '../../features/project/services';

import { MainLayout } from './main-layout';

describe('MainLayout', () => {
  let component: MainLayout;
  let fixture: ComponentFixture<MainLayout>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MainLayout],
      providers: [
        provideRouter([]),
        {
          provide: AuthService,
          useValue: { currentUser: signal(null), logout: vi.fn() },
        },
        {
          provide: ClientService,
          useValue: {
            state: signal({ clients: new Map(), loading: false }),
            stateWithProjects: signal({ clientsWithProjects: [], loading: false }),
            getAllClients: vi.fn(),
          },
        },
        {
          provide: projectService,
          useValue: {
            state: signal({ projects: new Map(), loading: false }),
            getAll: vi.fn(),
          },
        },
        {
          provide: TopbarTitleService,
          useValue: {
            projectContext: signal(null),
            setProjectBreadcrumb: vi.fn(),
            clearProjectBreadcrumb: vi.fn(),
          },
        },
        {
          provide: PermissionService,
          useValue: { load: vi.fn(() => of(false)), can: vi.fn(() => false) },
        },
        { provide: MatDialog, useValue: { open: vi.fn() } },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(MainLayout);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
