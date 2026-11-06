import { provideHttpClient } from '@angular/common/http';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { TestBed } from '@angular/core/testing';
import { vi } from 'vitest';

import { ClientService } from './client-service';

vi.mock('sweetalert2', () => {
  const instance = {
    fire: vi.fn().mockResolvedValue({}),
    getContainer: vi.fn(),
    mixin: vi.fn(),
  };
  instance.mixin.mockReturnValue(instance);

  return { default: instance };
});

describe('ClientService', () => {
  let service: ClientService;
  let httpMock: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(ClientService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should cancel a previous list request when a search starts', () => {
    service.getAllClients();
    const listRequest = httpMock.expectOne('http://localhost:4100/client');

    service.searchClients('acme');
    const searchRequest = httpMock.expectOne('http://localhost:4100/client/search?q=acme');

    expect(listRequest.cancelled).toBe(true);

    searchRequest.flush([
      {
        id: 1,
        name: 'Acme',
        created_at: '2026-01-01T00:00:00Z',
      },
    ]);

    expect(service.state().clients.get(1)?.name).toBe('Acme');
    expect(service.state().loading).toBe(false);
  });

  it('should keep a newly created client in the local state', () => {
    service.getAllClients();
    httpMock.expectOne('http://localhost:4100/client').flush([]);

    service.createClient({ name: 'Nuevo cliente' });
    httpMock.expectOne('http://localhost:4100/client/new').flush({
      id: 2,
      name: 'Nuevo cliente',
      created_at: '2026-01-01T00:00:00Z',
    });
    httpMock.expectOne('http://localhost:4100/client/with-projects').flush([]);

    expect(service.state().clients.get(2)?.name).toBe('Nuevo cliente');
  });
});
