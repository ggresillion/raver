import {Injectable} from '@angular/core';
import * as io from 'socket.io-client';
import {environment} from '../../environments/environment';
import {BotStatus} from '../shared/model/bot-status';
import {Observable, Subject} from 'rxjs';
import {Socket} from 'socket.io';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {

  private socket: Socket;
  private statusSubject = new Subject<BotStatus>();

  constructor() {
    this.socket = io(environment.websocket + 'bot');
    this.bindToEvents();
  }

  private bindToEvents() {
    this.socket.on('status', status => this.statusSubject.next(status));
  }

  public getStatus(): Observable<BotStatus> {
    return this.statusSubject.asObservable();
  }
}
