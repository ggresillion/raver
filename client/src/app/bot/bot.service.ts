import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable} from 'rxjs';
import {environment} from '../../environments/environment';
import {Server} from 'socket.io';
import {ServerEvents} from './model/server-events.enum';
import * as io from 'socket.io-client';
import {BotStateDTO} from './model/bot-state.dto';

@Injectable({
  providedIn: 'root'
})
export class BotService {

  private socket: Server;
  private stateSubject = new BehaviorSubject<BotStateDTO>(null);

  constructor() {
    this.socket = io(environment.websocket + 'player');
    this.bindToEvents();
  }

  public getState(): Observable<BotStateDTO> {
    return this.stateSubject.asObservable();
  }

  private bindToEvents(): void {
    this.socket.on(ServerEvents.STATE_UPDATE.toString(), data => this.onStateUpdate(data.state));
  }

  private onStateUpdate(state: BotStateDTO): void {
    this.stateSubject.next(state);
  }
}
