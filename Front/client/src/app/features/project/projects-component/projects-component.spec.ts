import { ComponentFixture, TestBed } from '@angular/core/testing';

import { projectsComponent } from './projects-component';

describe('projectsComponent', () => {
  let component: projectsComponent;
  let fixture: ComponentFixture<projectsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [projectsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(projectsComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
