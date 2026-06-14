import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WizardModalComponent } from './wizard-modal-component';

describe('WizardModalComponent', () => {
  let component: WizardModalComponent;
  let fixture: ComponentFixture<WizardModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WizardModalComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(WizardModalComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
