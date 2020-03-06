import {NgModule} from '@angular/core';
import {HomeComponent} from './home.component';
import {SharedModule} from '../shared/shared.module';
import {SoundModule} from '../sound/sound.module';
import {BotStatusComponent} from './components/bot-status/bot-status.component';
import {RouterModule} from '@angular/router';
import {GuildsModule} from '../guilds/guilds.module';
import {AddBotToGuildDialogComponent} from '../guilds/dialogs/add-bot-to-guild-dialog/add-bot-to-guild-dialog.component';
import {BotModule} from '../bot/bot.module';

@NgModule({
  declarations: [
    HomeComponent,
    BotStatusComponent,
    AddBotToGuildDialogComponent,
  ],
  entryComponents: [
    AddBotToGuildDialogComponent,
  ],
  imports: [
    SharedModule,
    SoundModule,
    RouterModule,
    GuildsModule,
    BotModule,
  ]
})
export class HomeModule {
}
