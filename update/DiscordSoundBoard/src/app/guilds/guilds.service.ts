import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Guild} from '../models/guild';
import {environment} from '../../environments/environment';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class GuildsService {

  constructor(private readonly http: HttpClient) {
  }

  public getAvailableGuilds(): Observable<Guild[]> {
    return this.http.get<Guild[]>(environment.api + '/guilds');
  }
}
