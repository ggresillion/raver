import { Component, OnInit } from '@angular/core';
import { YoutubeService } from '../../youtube.service';
import { TrackInfos } from '../../model/track-infos';

@Component({
  selector: 'app-search-youtube',
  templateUrl: './search-youtube.component.html',
  styleUrls: ['./search-youtube.component.scss']
})
export class SearchYoutubeComponent implements OnInit {

  public searchString = '';
  public videos: TrackInfos[] = [];
  public loading = false;

  constructor(
    private readonly youtubeService: YoutubeService,
  ) {
  }

  public ngOnInit(): void {
  }

  public search(): void {
    this.loading = true;
    this.youtubeService.searchYoutube(this.searchString)
      .subscribe(videos => {
        this.videos = videos;
        this.loading = false;
      });
  }

  public addToPlaylist(video: TrackInfos): void {
    this.youtubeService.addToPlaylist(video);
  }

}
