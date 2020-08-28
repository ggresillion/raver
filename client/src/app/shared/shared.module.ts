import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { DragDropModule } from '@angular/cdk/drag-drop';
import { DraggableDirective } from './directives/draggable.directive';
import { DropperDirective } from './directives/dropper.directive';
import { CategoryService } from './services/category.service';
import { BotService } from './services/bot.service';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatMenuModule } from '@angular/material/menu';
import { MatButtonModule } from '@angular/material/button';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatDialogModule } from '@angular/material/dialog';
import { MatInputModule } from '@angular/material/input';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatSlideToggleModule } from '@angular/material/slide-toggle';
import { MatExpansionModule } from '@angular/material/expansion';
import { MatSliderModule } from '@angular/material/slider';
import { MatTableModule } from '@angular/material/table';
import { MatPaginatorModule } from '@angular/material/paginator';
import { FlexLayoutModule } from '@angular/flex-layout';

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
    MatSliderModule,
    MatExpansionModule,
    MatTableModule,
    MatPaginatorModule,
    FlexLayoutModule,
    ReactiveFormsModule
  ],
  declarations: [DraggableDirective, DropperDirective],
  providers: [CategoryService, BotService],
})
export class SharedModule {
}
