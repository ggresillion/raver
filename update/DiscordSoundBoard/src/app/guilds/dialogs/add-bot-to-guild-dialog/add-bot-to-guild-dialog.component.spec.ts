import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddBotToGuildDialogComponent } from './add-bot-to-guild-dialog.component';

describe('AddBotToGuildDialogComponent', () => {
  let component: AddBotToGuildDialogComponent;
  let fixture: ComponentFixture<AddBotToGuildDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AddBotToGuildDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddBotToGuildDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
