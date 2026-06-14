import { TestBed } from '@angular/core/testing';

import { TypeFiles } from './type-files';

describe('TypeFiles', () => {
  let service: TypeFiles;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(TypeFiles);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
