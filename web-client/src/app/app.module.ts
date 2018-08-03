import {NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {HttpClientModule} from '@angular/common/http';
import { RemoveExtensionPipe } from './pipes/remove-extension/remove-extension.pipe';
import {
  MatButtonModule, MatCardModule,
  MatCommonModule,
  MatExpansionModule, MatFormFieldModule,
  MatIconModule, MatInputModule,
  MatMenuModule, MatProgressBarModule,
  MatSidenavModule, MatSnackBarModule, MatTabsModule,
  MatToolbarModule
} from "@angular/material";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {BrowserModule} from "@angular/platform-browser";

@NgModule({
  declarations: [
    AppComponent,
    RemoveExtensionPipe
  ],
  imports: [
    BrowserModule,
    MatToolbarModule,
    MatIconModule,
    MatCommonModule,
    MatMenuModule,
    MatSidenavModule,
    MatExpansionModule,
    MatFormFieldModule,
    MatButtonModule,
    MatTabsModule,
    MatCardModule,
    MatInputModule,
    BrowserAnimationsModule,
    MatProgressBarModule,
    MatSnackBarModule,
    HttpClientModule
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
