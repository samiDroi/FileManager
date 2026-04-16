import { TestBed } from '@angular/core/testing';

import { ProyectService } from './proyect-service';

describe('ProyectService', () => {
  let service: ProyectService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ProyectService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
