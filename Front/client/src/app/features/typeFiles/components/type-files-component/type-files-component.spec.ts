import { ComponentFixture, TestBed } from '@angular/core/testing';

import { TypeFilesComponent } from './type-files-component';

describe('TypeFilesComponent', () => {
  let component: TypeFilesComponent;
  let fixture: ComponentFixture<TypeFilesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TypeFilesComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(TypeFilesComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
