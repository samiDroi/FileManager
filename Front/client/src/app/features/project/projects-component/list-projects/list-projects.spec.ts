import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListProjects } from './list-projects';

describe('ListProjects', () => {
  let component: ListProjects;
  let fixture: ComponentFixture<ListProjects>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListProjects],
    }).compileComponents();

    fixture = TestBed.createComponent(ListProjects);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
