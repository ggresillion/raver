import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {HomeComponent} from './home.component';
import {SharedModule} from '../shared/shared.module';
import {SoundModule} from '../sound/sound.module';
import {BotStatusComponent} from './components/bot-status/bot-status.component';

@NgModule({
  declarations: [HomeComponent, BotStatusComponent],
  imports: [
    CommonModule,
    SharedModule,
    SoundModule,
  ]
})
export class HomeModule {
}
