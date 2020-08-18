import {Component, ElementRef, Input, OnInit, ViewChild} from '@angular/core';
import {YoutubeService} from '../../youtube.service';
import {TrackInfos} from '../../model/track-infos';
import {PlayerStatus} from '../../model/player-status';
import {MatTableDataSource} from '@angular/material/table';
import {MatPaginator} from '@angular/material/paginator';

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
  pageSizeOptions = [10, 20, 30];

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

  constructor(private playlistService: YoutubeService) {
  }

  ngOnInit() {
    this.setDataSourceAttributes();
    this.playlistService.getState().subscribe((state) => {
      this.status = state.status;
      this.playlist = state.playlist;
      if (state.totalLengthSeconds) {
        this.duration = Math.trunc(state.totalLengthSeconds);
      }
      this.setDataSourceAttributes();
    });
    this.playlistService.getProgress().subscribe(progress => {
      this.currentTime = Math.trunc(progress);
    });
  }

  @ViewChild(MatPaginator) set matPaginator(mp: MatPaginator) {
    this.paginator = mp;
    this.setDataSourceAttributes();
  }

  setDataSourceAttributes() {
    if (this.playlist) {
      const data = [...this.playlist];
      data.splice(0, 1);
      this.dataSource = new MatTableDataSource<TrackInfos>(data);
      this.dataSource.paginator = this.paginator;
    }
  }

  nextSong(): void {
    if (((this.playlistService.indexSong + 1) % this.paginator.pageSize) === 0 ||
      (this.playlistService.indexSong + 1) === this.paginator.length) {
      if (this.paginator.hasNextPage()) {
        this.paginator.nextPage();
      } else if (!this.paginator.hasNextPage()) {
        this.paginator.firstPage();
      }
    }
    this.currentTime = 0;
    this.duration = 0.01;
    this.playlistService.nextSong();
  }

  previousSong(): void {
    this.currentTime = 0;
    this.duration = 0.01;
    if (!this.checkIfSongHasStartedSinceAtleastTwoSeconds()) {
      if (((this.playlistService.indexSong) % this.paginator.pageSize) === 0 ||
        (this.playlistService.indexSong) === 0) {
        if (this.paginator.hasPreviousPage()) {
          this.paginator.previousPage();
        } else if (!this.paginator.hasPreviousPage()) {
          this.paginator.lastPage();
        }
      }
      this.playlistService.previousSong();
    } else {
    }
  }

  selectTrack(index: number): void {
    this.playlistService.selectATrack(index);
    setTimeout(() => {
      this.player.nativeElement.play();
    }, 0);
  }

  checkIfSongHasStartedSinceAtleastTwoSeconds(): boolean {
    return this.player.nativeElement.currentTime > 2;
  }

  currTimePosChanged(event) {
    this.player.nativeElement.currentTime = event.value;
  }

  bindPlayerEvent(): void {
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

  playPause() {
    if (this.status === PlayerStatus.PLAYING) {
      this.playlistService.pause();
    } else {
      this.playlistService.play();
    }
  }
}
