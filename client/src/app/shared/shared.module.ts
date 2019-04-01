import {NgModule} from '@angular/core';
import {MatCardModule, MatIconModule, MatTabsModule, MatToolbarModule} from '@angular/material';
import {CommonModule} from '@angular/common';

@NgModule({
  exports: [
    CommonModule,
    FormsModule,
    MatCardModule,
    MatIconModule,
    MatTabsModule,
    MatCardModule,
    MatToolbarModule,
  ]
})
export class SharedModule {
}
