import {NgModule} from '@angular/core';
import {YoutubeComponent} from './youtube.component';
import {SharedModule} from '../shared/shared.module';
import {SecondsToMinutesPipe} from './seconds-to-minutes.pipe';
import {MusicPlayerComponent} from './components/music-player/music-player.component';
import {YoutubeService} from './youtube.service';
import { SearchYoutubeComponent } from './components/search-youtube/search-youtube.component';
import { YoutubeThumbnailComponent } from './components/youtube-thumbnail/youtube-thumbnail.component';

@NgModule({
  declarations: [YoutubeComponent, SecondsToMinutesPipe, MusicPlayerComponent, SearchYoutubeComponent, YoutubeThumbnailComponent],
  imports: [
    SharedModule
  ],
  providers: [YoutubeService]
})
export class YoutubeModule {
}
