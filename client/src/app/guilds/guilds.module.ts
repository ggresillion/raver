import {NgModule} from '@angular/core';
import {GuildsComponent} from './guilds.component';
import {SharedModule} from '../shared/shared.module';

@NgModule({
  declarations: [GuildsComponent],
  imports: [
    SharedModule,
  ]
})
export class GuildsModule {
}
