import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable} from 'rxjs';
import * as io from 'socket.io-client';
import {environment} from '../../environments/environment';
import {PlayerStatus} from './model/player-status';
import {HttpClient} from '@angular/common/http';
import {ClientEvents} from './model/client-events.enum';
import {ServerEvents} from './model/server-events.enum';
import {Socket} from 'socket.io';
import {TrackInfos} from './model/track-infos';
import {PlayerState} from './model/player-state';

@Injectable({
  providedIn: 'root'
})
export class YoutubeService {

  private playlistSubject: BehaviorSubject<TrackInfos[]> = new BehaviorSubject([]);
  public indexSong = 0;
  private statusSubject: BehaviorSubject<PlayerStatus> = new BehaviorSubject(PlayerStatus.IDLE);
  private currentTime = 0;
  private duration = 0;
  private socket: Socket;

  constructor(
    private http: HttpClient,
  ) {
    this.socket = io(environment.websocket + 'player');
    this.bindToEvents();
  }

  public addToPlaylist(track: TrackInfos) {
    this.socket.emit(ClientEvents.ADD_TO_PLAYLIST, {track});
  }

  public getPlaylist() {
    return this.playlistSubject.asObservable();
  }

  public nextSong(): void {
    if ((this.indexSong + 1) >= this.playlistSubject.value.length) {
      this.indexSong = 0;
    } else {
      this.indexSong++;
    }
    this.socket.emit(ClientEvents.NEXT);
  }

  public previousSong(): void {
    if ((this.indexSong - 1) < 0) {
      this.indexSong = (this.playlistSubject.value.length - 1);
    } else {
      this.indexSong--;
    }
  }

  public resetPlaylist(): void {
    this.indexSong = 0;
  }

  public selectATrack(index: number): void {
    this.indexSong = index - 1;
  }

  public getStatus() {
    return this.statusSubject.asObservable();
  }

  public searchYoutube(searchString: string): Observable<any[]> {
    return this.http.get<any[]>(`${environment.api}/youtube/search?q=${searchString}`);
  }

  public play() {
    this.socket.emit(ClientEvents.PLAY);
  }

  public pause() {
    this.socket.emit(ClientEvents.PAUSE);
  }

  private bindToEvents() {
    this.socket.on(ServerEvents.SYNC, infos => {
      this.statusSubject.next(infos.status);
      this.playlistSubject.next(infos.playlist);
    });
    this.socket.on(ServerEvents.STATUS_UPDATED, data => this.onStatusUpdate(data.status));
    this.socket.on(ServerEvents.ADD_TO_PLAYLIST, data => this.onAddToPlaylist(data.track));
    this.socket.on(ServerEvents.GET_PLAYLIST, playlist => this.playlistSubject.next(playlist));
    this.socket.on(ServerEvents.STATE_UPDATED, data => this.onStateUpdate(data.state));
  }

  private onStatusUpdate(status: PlayerStatus) {
    this.statusSubject.next(status);
  }

  private onStateUpdate(state: PlayerState) {
    console.log('stateUpdated:', state);
    this.onStatusUpdate(state.status);
    this.playlistSubject.next(state.playlist);
  }

  private onAddToPlaylist(track) {
    this.playlistSubject.next([...this.playlistSubject.value, track]);
  }
}
