import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {LoginComponent} from './login/login.component';
import {HomeComponent} from './home/home.component';
import {IsConnectedGuard} from './auth/is-connected.guard';
import {SoundComponent} from './sound/sound.component';
import {YoutubeComponent} from './youtube/youtube.component';

const routes: Routes = [
  {path: 'login', component: LoginComponent},
  {
    path: '', component: HomeComponent, canActivate: [IsConnectedGuard], children: [
      {path: '', component: SoundComponent},
      {path: 'youtube', component: YoutubeComponent},
    ]
  }
];

@NgModule({
  imports: [RouterModule.forRoot(routes, {useHash: true})],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
