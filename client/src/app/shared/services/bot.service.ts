import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../../environments/environment';
import {Observable} from 'rxjs';
import {BotInfos} from '../model/bot-infos';
import {BotStatus} from '../model/bot-status';
import {WebsocketService} from '../../websocket/websocket.service';

@Injectable()
export class BotService {

  constructor(
    private http: HttpClient,
    private readonly websocketService: WebsocketService) {
  }

  public getInfos(): Observable<BotInfos> {
    return this.http.get<BotInfos>(`${environment.api}/bot/infos`);
  }

  public getStatus(): Observable<BotStatus> {
    return this.websocketService.getStatus();
  }
}
