import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../../environments/environment';
import {Observable} from 'rxjs';
import {BotInfos} from '../model/bot-infos';

@Injectable()
export class BotService {

  constructor(private http: HttpClient) {
  }

  public getInfos(): Observable<BotInfos> {
    return this.http.get<BotInfos>(`${environment.api}/bot/infos`);
  }
}
