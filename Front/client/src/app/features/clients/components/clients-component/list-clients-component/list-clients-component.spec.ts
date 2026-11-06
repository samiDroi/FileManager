import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListClientsComponent } from './list-clients-component';

describe('ListClientsComponent', () => {
  let component: ListClientsComponent;
  let fixture: ComponentFixture<ListClientsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ListClientsComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(ListClientsComponent);
    fixture.componentRef.setInput('clients', []);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
