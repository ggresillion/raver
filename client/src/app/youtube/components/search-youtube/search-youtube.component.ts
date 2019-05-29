import { Component, OnInit } from '@angular/core';
import { YoutubeService } from '../../youtube.service';
import { TrackInfos } from '../../model/track-infos';

@Component({
  selector: 'app-search-youtube',
  templateUrl: './search-youtube.component.html',
  styleUrls: ['./search-youtube.component.sass']
})
export class SearchYoutubeComponent implements OnInit {

  searchString = '';
  videos = [];

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
