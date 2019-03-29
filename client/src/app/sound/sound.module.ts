import {NgModule} from '@angular/core';
import {SoundService} from './sound.service';
import {SoundComponent} from './sound.component';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [SoundComponent],
  imports: [
    SharedModule,
  ],
  providers: [SoundService],
  exports: [SoundComponent]
})
export class SoundModule {
}
