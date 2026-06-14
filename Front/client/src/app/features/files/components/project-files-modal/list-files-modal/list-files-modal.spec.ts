import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListFilesModal } from './list-files-modal';

describe('ListFilesModal', () => {
  let component: ListFilesModal;
  let fixture: ComponentFixture<ListFilesModal>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListFilesModal],
    }).compileComponents();

    fixture = TestBed.createComponent(ListFilesModal);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
