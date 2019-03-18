import {NgModule} from '@angular/core';

import {AppComponent} from './app.component';
import {HttpClientModule} from '@angular/common/http';
import {RemoveExtensionPipe} from './pipes/remove-extension/remove-extension.pipe';
import {
  MatButtonModule,
  MatCardModule,
  MatCommonModule,
  MatDialogModule,
  MatExpansionModule,
  MatFormFieldModule,
  MatIconModule,
  MatInputModule,
  MatMenuModule,
  MatProgressBarModule,
  MatProgressSpinnerModule,
  MatSidenavModule,
  MatSnackBarModule,
  MatTabsModule,
  MatToolbarModule
} from '@angular/material';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {BrowserModule} from '@angular/platform-browser';
import {AddFromYoutubeDialogComponent} from './dialogs/add-from-youtube-dialog/add-from-youtube-dialog.component';
import {FormsModule} from '@angular/forms';
import {CreateCategoryDialogComponent} from './dialogs/create-category-dialog/create-category-dialog.component';
import {LoginComponent} from './components/login/login.component';
import {AppRoutingModule} from './app-routing.module';
import {HomeComponent} from './components/home/home.component';
import { SongsComponent } from './components/home/songs/songs.component';
import { GuildsComponent } from './components/home/guilds/guilds.component';

@NgModule({
  declarations: [
    AppComponent,
    RemoveExtensionPipe,
    AddFromYoutubeDialogComponent,
    CreateCategoryDialogComponent,
    LoginComponent,
    HomeComponent,
    SongsComponent,
    GuildsComponent
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
    HttpClientModule,
    MatDialogModule,
    FormsModule,
    MatProgressSpinnerModule,
    AppRoutingModule
  ],
  providers: [],
  bootstrap: [AppComponent],
  entryComponents: [
    AddFromYoutubeDialogComponent,
    CreateCategoryDialogComponent
  ]
})
export class AppModule {
}
