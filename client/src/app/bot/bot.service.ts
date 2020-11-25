import { Injectable } from '@angular/core';
import { BehaviorSubject, Observable, ReplaySubject } from 'rxjs';
import { environment } from '../../environments/environment';
import { Server } from 'socket.io';
import { ServerEvents } from './model/server-events.enum';
import * as io from 'socket.io-client';
import { BotStateDTO } from './model/bot-state.dto';
import { HttpClient } from '@angular/common/http';
import { ClientEvents } from './model/client-events.enum';
import { GuildsService } from '../guilds/guilds.service';
import { Guild } from '../models/guild';

@Injectable({
  providedIn: 'root'
})
export class BotService {

  private socket: Server;
  private stateSubject = new ReplaySubject<BotStateDTO>(1);
  private currentGuild: Guild;

  constructor(
    private readonly http: HttpClient,
    private readonly guildsService: GuildsService) {
    this.socket = io(environment.websocket + 'bot')
      .on('connect', () => {
        this.guildsService.getSelectedGuild().subscribe(g => {
          this.currentGuild = g;
          this.socket.emit(ServerEvents.JOIN_ROOM, { guildId: g.id });
        });
        this.bindToEvents();
      });
  }

  public getState(): Observable<BotStateDTO> {
    return this.stateSubject.asObservable();
  }

  public joinMyChannel(): Observable<void> {
    return this.http.post<void>(environment.api + '/bot/join', null);
  }

  public setVolume(volume: number) {
    this.socket.emit(ClientEvents.CHANGE_VOLUME.toString(), { volume });
  }

  private bindToEvents(): void {
    this.socket.on(ServerEvents.STATE_UPDATE.toString(), data => this.onStateUpdate(data.state));
  }

  private onStateUpdate(state: BotStateDTO): void {
    this.stateSubject.next(state);
  }
}
