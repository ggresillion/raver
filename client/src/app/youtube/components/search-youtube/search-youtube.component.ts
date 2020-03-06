import { Component, OnInit } from '@angular/core';
import { YoutubeService } from '../../youtube.service';
import { TrackInfos } from '../../model/track-infos';

@Component({
  selector: 'app-search-youtube',
  templateUrl: './search-youtube.component.html',
  styleUrls: ['./search-youtube.component.scss']
})
export class SearchYoutubeComponent implements OnInit {

  searchString = '';
  videos: TrackInfos[] = [];

  constructor(
    private readonly youtubeService: YoutubeService,
  ) {
  }

  ngOnInit() {
  }

  search() {
    this.youtubeService.searchYoutube(this.searchString)
      .subscribe(videos => this.videos = videos);
  }

  addToPlaylist(video: TrackInfos) {
    this.youtubeService.addToPlaylist(video);
  }

}
