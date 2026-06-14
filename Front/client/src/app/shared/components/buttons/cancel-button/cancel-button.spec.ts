import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CancelButton } from './cancel-button';

describe('CancelButton', () => {
  let component: CancelButton;
  let fixture: ComponentFixture<CancelButton>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CancelButton],
    }).compileComponents();

    fixture = TestBed.createComponent(CancelButton);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
