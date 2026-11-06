import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, Router } from '@angular/router';
import { of } from 'rxjs';
import { AuthService } from '../../auth-service';

import { RecoveryPasswordComponent } from './recovery-password-component';

describe('RecoveryPasswordComponent', () => {
  let component: RecoveryPasswordComponent;
  let fixture: ComponentFixture<RecoveryPasswordComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RecoveryPasswordComponent],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: { queryParams: of({ token: 'test-token' }) },
        },
        { provide: Router, useValue: { navigate: vi.fn() } },
        { provide: AuthService, useValue: { recoveryPassword: vi.fn() } },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(RecoveryPasswordComponent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
