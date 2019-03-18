import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddFromYoutubeDialogComponent } from './add-from-youtube-dialog.component';

describe('AddFromYoutubeDialogComponent', () => {
  let component: AddFromYoutubeDialogComponent;
  let fixture: ComponentFixture<AddFromYoutubeDialogComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AddFromYoutubeDialogComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddFromYoutubeDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
