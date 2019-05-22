import {Component, OnInit, Input, ViewChild, ElementRef} from '@angular/core';
import {MatSlider, MatPaginator, MatTableDataSource} from '@angular/material';
import {YoutubeService} from '../../youtube.service';
import {Status} from '../../model/status';

@Component({
  selector: 'app-music-player',
  templateUrl: './music-player.component.html',
  styleUrls: ['./music-player.component.scss']
})
export class MusicPlayerComponent implements OnInit {

  @Input()
  displayTitle = false;

  @Input()
  displayPlaylist = false;

  @Input()
  pageSizeOptions = [10, 20, 30];

  @Input()
  expanded = true;

  @ViewChild('audioPlayer') player: ElementRef;

  displayedColumns: string[] = ['title', 'status'];

  dataSource = new MatTableDataSource<any>();

  paginator: MatPaginator;

  playlistData: any[];

  playlistTrack: any;
  loaderDisplay = false;
  isPlaying = false;
  currentTime = 0;
  duration = 0.01;

  public status = {status: Status.IDLE};

  constructor(private playlistService: YoutubeService) {
    this.playlistService.getStatus().subscribe((status) => {
      this.status = status;
    });
  }

  ngOnInit() {
    this.setDataSourceAttributes();
    this.bindPlayerEvent();
    this.player.nativeElement.addEventListener('ended', () => {
      if (this.checkIfSongHasStartedSinceAtleastTwoSeconds()) {
        this.nextSong();
      }
    });
    this.playlistService.getSubjectCurrentTrack().subscribe((playlistTrack) => {
      this.playlistTrack = playlistTrack;
    });
    this.player.nativeElement.currentTime = 0;
    this.playlistService.init();
  }

  @ViewChild(MatPaginator) set matPaginator(mp: MatPaginator) {
    this.paginator = mp;
    this.setDataSourceAttributes();
  }

  setDataSourceAttributes() {
    let index = 1;
    if (this.playlistData) {
      this.playlistData.forEach(data => {
        data.index = index++;
      });
      this.dataSource = new MatTableDataSource<any>(this.playlistData);
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
    this.play();
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
      this.resetSong();
    }
    this.play();
  }

  resetSong(): void {
    this.player.nativeElement.src = this.playlistTrack[1].link;
  }

  selectTrack(index: number): void {
    console.log('selectTrack(index: number): void: ' + index);
    this.playlistService.selectATrack(index);
    setTimeout(() => {
      this.player.nativeElement.play();
    }, 0);
  }

  checkIfSongHasStartedSinceAtleastTwoSeconds(): boolean {
    return this.player.nativeElement.currentTime > 2;
  }

  @Input()
  set playlist(playlist: any[]) {
    this.playlistData = playlist;
    this.ngOnInit();
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

  playBtnHandler(): void {
    if (this.loaderDisplay) {
      return;
    }
    if (this.player.nativeElement.paused) {
      this.player.nativeElement.play(this.currentTime);
    } else {
      this.currentTime = this.player.nativeElement.currentTime;
      this.player.nativeElement.pause();
    }
  }

  play(): void {
    setTimeout(() => {
      this.player.nativeElement.play();
    }, 0);
  }
}
