import {NgModule} from '@angular/core';
import {SoundService} from './sound.service';
import {SoundComponent} from './sound.component';
import {SharedModule} from '../shared/shared.module';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';
import {AddFromYoutubeDialogComponent} from './dialogs/add-from-youtube-dialog/add-from-youtube-dialog.component';
import {CreateCategoryDialogComponent} from './dialogs/create-category-dialog/create-category-dialog.component';

@NgModule({
  declarations: [
    SoundComponent,
    UploadSoundDialogComponent,
    AddFromYoutubeDialogComponent,
    CreateCategoryDialogComponent,
  ],
  imports: [
    SharedModule,
  ],
  providers: [SoundService],
  exports: [SoundComponent],
  entryComponents: [
    UploadSoundDialogComponent,
    AddFromYoutubeDialogComponent,
    CreateCategoryDialogComponent,
  ],
})
export class SoundModule {
}
