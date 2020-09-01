import { Component, ElementRef, Input, OnInit, ViewChild } from '@angular/core';
import { YoutubeService } from '../../youtube.service';
import { TrackInfos } from '../../model/track-infos';
import { PlayerStatus } from '../../model/player-status';
import { MatTableDataSource } from '@angular/material/table';
import { MatPaginator } from '@angular/material/paginator';

@Component({
  selector: 'app-music-player',
  templateUrl: './music-player.component.html',
  styleUrls: ['./music-player.component.scss']
})
export class MusicPlayerComponent implements OnInit {

  @Input()
  displayTitle = true;

  @Input()
  displayPlaylist = true;

  @Input()
  pageSizeOptions = [4];

  @Input()
  expanded = true;

  @ViewChild('audioPlayer') player: ElementRef;

  displayedColumns: string[] = ['position', 'thumbnail', 'title'];

  dataSource = new MatTableDataSource<TrackInfos>();

  paginator: MatPaginator;

  playlist: TrackInfos[];
  loaderDisplay = false;
  isPlaying = false;
  currentTime = 0;
  duration;

  public status = PlayerStatus.IDLE;

  constructor(private youtubeService: YoutubeService) {
  }

  public ngOnInit() {
    this.setDataSourceAttributes();
    this.youtubeService.getState().subscribe((state) => {
      this.status = state.status;
      this.playlist = state.playlist;
      if (state.totalLengthSeconds) {
        this.duration = Math.trunc(state.totalLengthSeconds);
      }
      this.setDataSourceAttributes();
    });
    this.youtubeService.getProgress().subscribe(progress => {
      this.currentTime = Math.trunc(progress);
    });
  }

  @ViewChild(MatPaginator) set matPaginator(mp: MatPaginator) {
    this.paginator = mp;
    this.setDataSourceAttributes();
  }

  public setDataSourceAttributes() {
    if (this.playlist) {
      const data = [...this.playlist];
      data.splice(0, 1);
      this.dataSource = new MatTableDataSource<TrackInfos>(data);
      this.dataSource.paginator = this.paginator;
    }
  }

  public nextSong(): void {
    if (((this.youtubeService.indexSong + 1) % this.paginator.pageSize) === 0 ||
      (this.youtubeService.indexSong + 1) === this.paginator.length) {
      if (this.paginator.hasNextPage()) {
        this.paginator.nextPage();
      } else if (!this.paginator.hasNextPage()) {
        this.paginator.firstPage();
      }
    }
    this.currentTime = 0;
    this.duration = 0.01;
    this.youtubeService.nextSong();
  }

  public previousSong(): void {
    this.currentTime = 0;
    this.duration = 0.01;
    if (!this.checkIfSongHasStartedSinceAtleastTwoSeconds()) {
      if (((this.youtubeService.indexSong) % this.paginator.pageSize) === 0 ||
        (this.youtubeService.indexSong) === 0) {
        if (this.paginator.hasPreviousPage()) {
          this.paginator.previousPage();
        } else if (!this.paginator.hasPreviousPage()) {
          this.paginator.lastPage();
        }
      }
      this.youtubeService.previousSong();
    } else {
    }
  }

  public selectTrack(index: number): void {
    this.youtubeService.selectATrack(index);
    setTimeout(() => {
      this.player.nativeElement.play();
    }, 0);
  }

  public checkIfSongHasStartedSinceAtleastTwoSeconds(): boolean {
    return this.player.nativeElement.currentTime > 2;
  }

  public currTimePosChanged(event) {
    this.player.nativeElement.currentTime = event.value;
  }

  public bindPlayerEvent(): void {
    this.player.nativeElement.addEventListener('playing', () => {
      this.isPlaying = true;
      this.duration = Math.floor(this.player.nativeElement.duration);
    });
    this.player.nativeElement.addEventListener('pause', () => {
      this.isPlaying = false;
    });
    this.player.nativeElement.addEventListener('timeupdate', () => {
      this.currentTime = Math.floor(this.player.nativeElement.currentTime);
    });
    this.player.nativeElement.addEventListener('loadstart', () => {
      this.loaderDisplay = true;
    });
    this.player.nativeElement.addEventListener('loadeddata', () => {
      this.loaderDisplay = false;
      this.duration = Math.floor(this.player.nativeElement.duration);
    });
  }

  public playPause() {
    if (this.status === PlayerStatus.PLAYING) {
      this.youtubeService.pause();
    } else {
      this.youtubeService.play();
    }
  }

  public stop() {
    this.youtubeService.stop();
  }

  public getTrackPosition(track: TrackInfos): number {
    return this.playlist.findIndex(t => t === track) + 1;
  }

  public moveUpwards(element: TrackInfos): void {
    this.youtubeService.moveTrackUpwards(this.getIndex(element));
  }

  public moveDownwards(element: TrackInfos): void {
    this.youtubeService.moveTrackDownwards(this.getIndex(element));
  }

  public removeFromPlaylist(element: TrackInfos): void {
    this.youtubeService.removeFromPlaylist(this.getIndex(element));
  }

  public getIndex(element: TrackInfos): number {
    return this.playlist.findIndex(t => t === element)
  }
}
