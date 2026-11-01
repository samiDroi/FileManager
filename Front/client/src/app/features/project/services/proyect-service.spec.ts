import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { vi } from 'vitest';

import { projectService } from './project-service';

vi.mock('sweetalert2', () => {
  const instance = {
    fire: vi.fn().mockResolvedValue({}),
    getContainer: vi.fn(),
    mixin: vi.fn(),
  };
  instance.mixin.mockReturnValue(instance);

  return { default: instance };
});

describe('projectService', () => {
  let service: projectService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(projectService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should cancel a previous list request when a search starts', () => {
    service.getAll();
    const listRequest = httpMock.expectOne('http://localhost:4100/project');

    service.searchProjects('acme');
    const searchRequest = httpMock.expectOne('http://localhost:4100/project/search?q=acme');

    expect(listRequest.cancelled).toBe(true);

    searchRequest.flush([
      {
        id: 1,
        name: 'Acme',
        client_id: 1,
        client_name: 'Cliente',
        count_files: 0,
        created_at: '2026-01-01T00:00:00Z',
      },
    ]);

    expect(service.state().projects.get(1)?.name).toBe('Acme');
    expect(service.state().loading).toBe(false);
  });

  it('should keep a newly created project in the local state', () => {
    service.getAll();
    httpMock.expectOne('http://localhost:4100/project').flush([]);

    service.onCreate({ name: 'Nuevo proyecto', client_id: 1 });
    httpMock.expectOne('http://localhost:4100/project/new').flush({
      id: 3,
      name: 'Nuevo proyecto',
      client_id: 1,
      client_name: 'Cliente',
      count_files: 0,
      created_at: '2026-01-01T00:00:00Z',
    });

    expect(service.state().projects.get(3)?.name).toBe('Nuevo proyecto');
  });
});
