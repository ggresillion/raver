import {NgModule} from '@angular/core';
import {YoutubeComponent} from './youtube.component';
import {SharedModule} from '../shared/shared.module';
import {MatExpansionModule, MatIconModule, MatPaginatorModule, MatSliderModule, MatTableModule} from '@angular/material';
import {SecondsToMinutesPipe} from './seconds-to-minutes.pipe';
import {MusicPlayerComponent} from './components/music-player/music-player.component';
import {YoutubeService} from './youtube.service';
import { SearchYoutubeComponent } from './components/search-youtube/search-youtube.component';

@NgModule({
  declarations: [YoutubeComponent, SecondsToMinutesPipe, MusicPlayerComponent, SearchYoutubeComponent],
  imports: [
    SharedModule,
    MatSliderModule,
    MatIconModule,
    MatExpansionModule,
    MatTableModule,
    MatPaginatorModule,
  ],
  providers: [YoutubeService]
})
export class YoutubeModule {
}
