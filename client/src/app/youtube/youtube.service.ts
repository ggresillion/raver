import {Injectable} from '@angular/core';
import {BehaviorSubject} from 'rxjs';
import * as io from 'socket.io-client';
import {environment} from '../../environments/environment';
import Socket = SocketIOClient.Socket;
import {PlayerStatus} from './model/player-status';
import {Status} from './model/status';

@Injectable({
  providedIn: 'root'
})
export class YoutubeService {

  playlist: any[] = [];

  indexSong = 0;
  currentTrackSubject: BehaviorSubject<{}> = new BehaviorSubject(this.playlist[this.indexSong]);
  statusSubject: BehaviorSubject<PlayerStatus> = new BehaviorSubject({status: Status.IDLE});
  currentTime = 0;
  duration = 0;
  socket: Socket;

  constructor() {
    this.socket = io(environment.websocket + 'player');
    this.bindToEvents();
  }

  init(): void {
    this.updateCurrentSong();
  }

  nextSong(): void {
    if ((this.indexSong + 1) >= this.playlist.length) {
      this.indexSong = 0;
    } else {
      this.indexSong++;
    }
    this.updateCurrentSong();
  }

  previousSong(): void {
    if ((this.indexSong - 1) < 0) {
      this.indexSong = (this.playlist.length - 1);
    } else {
      this.indexSong--;
    }
    this.updateCurrentSong();
  }

  resetPlaylist(): void {
    this.indexSong = 0;
    this.updateCurrentSong();
  }

  selectATrack(index: number): void {
    this.indexSong = index - 1;
    this.updateCurrentSong();
  }

  updateCurrentSong(): void {
    const current = this.playlist[this.indexSong];
    const previous = ((this.indexSong - 1) >= 0) ? this.playlist[this.indexSong - 1] : this.playlist[this.playlist.length - 1];
    const next = ((this.indexSong + 1) >= this.playlist.length) ? this.playlist[0] : this.playlist[this.indexSong + 1];

    this.currentTrackSubject.next([
      previous,
      current,
      next
    ]);
  }

  getSubjectCurrentTrack(): BehaviorSubject<{}> {
    return this.currentTrackSubject;
  }

  private bindToEvents() {
    this.socket.on('status', status => this.updateStatus(status));
  }

  private updateStatus(status: PlayerStatus) {
    this.statusSubject.next(status);
  }

  public getStatus() {
    return this.statusSubject.asObservable();
  }
}
