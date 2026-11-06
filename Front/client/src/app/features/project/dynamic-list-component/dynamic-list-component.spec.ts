import { TestBed } from '@angular/core/testing';
import { MatDialog } from '@angular/material/dialog';
import { provideRouter } from '@angular/router';
import { RouterTestingHarness } from '@angular/router/testing';
import { signal } from '@angular/core';
import { ProjectData } from '../../../core';
import { ClientService } from '../../../features/clients/services';
import { projectService } from '../services';
import { DynamicListComponent } from './dynamic-list-component';

describe('DynamicListComponent', () => {
  let projectState: ReturnType<typeof signal<{ projects: Map<number, ProjectData>; loading: boolean }>>;
  let getAllProjects: ReturnType<typeof vi.fn>;

  beforeEach(async () => {
    projectState = signal({ projects: new Map<number, ProjectData>(), loading: false });
    getAllProjects = vi.fn();

    await TestBed.configureTestingModule({
      imports: [DynamicListComponent],
      providers: [
        provideRouter([{ path: 'overview/client/:id', component: DynamicListComponent }]),
        {
          provide: projectService,
          useValue: {
            state: projectState,
            getAll: getAllProjects,
            onCreate: vi.fn(),
            onUpdate: vi.fn(),
          },
        },
        {
          provide: ClientService,
          useValue: {
            stateWithProjects: signal({
              clientsWithProjects: [{ id: 42, name: 'Cliente de prueba' }],
              loading: false,
            }),
            getClientsWithProjects: vi.fn(),
          },
        },
        { provide: MatDialog, useValue: { open: vi.fn() } },
      ],
    }).compileComponents();
  });

  it('loads on direct client route and reacts to newly added projects', async () => {
    const harness = await RouterTestingHarness.create();
    const component = await harness.navigateByUrl('/overview/client/42', DynamicListComponent);

    expect(component.clientId()).toBe(42);
    expect(getAllProjects).toHaveBeenCalledOnce();
    expect(harness.routeNativeElement?.querySelector('app-empty-state')).not.toBeNull();

    projectState.set({
      projects: new Map([[7, {
        id: 7,
        name: 'Proyecto nuevo',
        client_id: 42,
        client_name: 'Cliente de prueba',
        count_files: 0,
        created_at: '2026-01-01',
      }]]),
      loading: false,
    });
    harness.detectChanges();

    expect(component.projects().map(project => project.name)).toEqual(['Proyecto nuevo']);
    expect(harness.routeNativeElement?.querySelector('app-empty-state')).toBeNull();
  });
});
