import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {FormsModule} from '@angular/forms';
import {
  MatButtonModule,
  MatCardModule,
  MatDialogModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatMenuModule,
  MatProgressBarModule,
  MatProgressSpinnerModule,
  MatSidenavModule,
  MatSliderModule,
  MatSlideToggleModule,
  MatTabsModule,
  MatToolbarModule
} from '@angular/material';
import {DragDropModule} from '@angular/cdk/drag-drop';
import { DraggableDirective } from './directives/draggable.directive';
import { DropperDirective } from './directives/dropper.directive';

@NgModule({
  exports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatIconModule,
    MatTabsModule,
    MatCardModule,
    MatToolbarModule,
    MatButtonModule,
    MatProgressBarModule,
    MatMenuModule,
    MatProgressSpinnerModule,
    MatFormFieldModule,
    MatDialogModule,
    MatInputModule,
    MatSidenavModule,
    DragDropModule,
    DraggableDirective,
    DropperDirective,
    MatSlideToggleModule,
  ],
  declarations: [DraggableDirective, DropperDirective]
})
export class SharedModule {
}
