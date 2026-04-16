import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ProyectsComponent } from './proyects-component';

describe('ProyectsComponent', () => {
  let component: ProyectsComponent;
  let fixture: ComponentFixture<ProyectsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ProyectsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ProyectsComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
