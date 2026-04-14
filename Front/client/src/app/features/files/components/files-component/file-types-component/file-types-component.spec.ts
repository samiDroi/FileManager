import { ComponentFixture, TestBed } from '@angular/core/testing';

import { FileTypesComponent } from './file-types-component';

describe('FileTypesComponent', () => {
  let component: FileTypesComponent;
  let fixture: ComponentFixture<FileTypesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [FileTypesComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(FileTypesComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
