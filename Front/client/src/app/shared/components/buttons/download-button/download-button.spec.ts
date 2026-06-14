import { ComponentFixture, TestBed } from '@angular/core/testing';

import { DownloadButton } from './download-button';

describe('DownloadButton', () => {
  let component: DownloadButton;
  let fixture: ComponentFixture<DownloadButton>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [DownloadButton],
    }).compileComponents();

    fixture = TestBed.createComponent(DownloadButton);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
