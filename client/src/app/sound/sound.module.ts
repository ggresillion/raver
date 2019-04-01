import {NgModule} from '@angular/core';
import {SoundService} from './sound.service';
import {SoundComponent} from './sound.component';
import {SharedModule} from '../shared/shared.module';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';

@NgModule({
  declarations: [SoundComponent, UploadSoundDialogComponent],
  imports: [
    SharedModule,
  ],
  providers: [SoundService],
  exports: [SoundComponent],
  entryComponents: [UploadSoundDialogComponent],
})
export class SoundModule {
}
