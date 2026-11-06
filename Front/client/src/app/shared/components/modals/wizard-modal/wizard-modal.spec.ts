import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';

import { WizardModal } from './wizard-modal';

describe('WizardModal', () => {
  let component: WizardModal;
  let fixture: ComponentFixture<WizardModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [WizardModal],
      providers: [
        { provide: MAT_DIALOG_DATA, useValue: { steps: [] } },
        { provide: MatDialogRef, useValue: { close: vi.fn() } },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(WizardModal);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
