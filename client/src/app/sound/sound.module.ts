import {NgModule} from '@angular/core';
import {SoundService} from './sound.service';
import {SoundComponent} from './sound.component';
import {SharedModule} from '../shared/shared.module';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';
import {AddFromYoutubeDialogComponent} from './dialogs/add-from-youtube-dialog/add-from-youtube-dialog.component';
import {CreateCategoryDialogComponent} from './dialogs/create-category-dialog/create-category-dialog.component';
import {RenameCategoryDialogComponent} from './dialogs/rename-category-dialog/rename-category-dialog.component';
import {YoutubeModule} from '../youtube/youtube.module';
import { EditSoundDialogComponent } from './dialogs/edit-sound-dialog/edit-sound-dialog.component';

@NgModule({
  declarations: [
    SoundComponent,
    UploadSoundDialogComponent,
    AddFromYoutubeDialogComponent,
    CreateCategoryDialogComponent,
    RenameCategoryDialogComponent,
    EditSoundDialogComponent
  ],
  imports: [
    SharedModule,
    YoutubeModule,
  ],
  providers: [SoundService],
  exports: [SoundComponent],
  entryComponents: [
    UploadSoundDialogComponent,
    AddFromYoutubeDialogComponent,
    CreateCategoryDialogComponent,
    RenameCategoryDialogComponent,
    EditSoundDialogComponent
  ],
})
export class SoundModule {
}
