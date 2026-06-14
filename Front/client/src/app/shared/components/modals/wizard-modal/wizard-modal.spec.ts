import { ComponentFixture, TestBed } from '@angular/core/testing';

import { WizardModal } from './wizard-modal';

describe('WizardModal', () => {
  let component: WizardModal;
  let fixture: ComponentFixture<WizardModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WizardModal],
    }).compileComponents();

    fixture = TestBed.createComponent(WizardModal);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
