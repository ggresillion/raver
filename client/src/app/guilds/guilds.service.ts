import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Guild } from '../models/guild';
import { environment } from '../../environments/environment';
import { BehaviorSubject, Observable, EMPTY, ReplaySubject } from 'rxjs';
import { tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class GuildsService {

  private guildSubject = new ReplaySubject<Guild>();

  constructor(private readonly http: HttpClient) {
  }

  public getAvailableGuilds(): Observable<Guild[]> {
    return this.http.get<Guild[]>(environment.api + '/guilds')
      .pipe(tap(guilds => {
        guilds.sort(g1 => -g1.isBotInGuild);
      }))
      .pipe(tap(guilds => this.setSelectedGuild(guilds[0])));
  }

  public getSelectedGuild(): Observable<Guild> {
    return this.guildSubject.asObservable();
  }

  public setSelectedGuild(guild: Guild): void {
    this.guildSubject.next(guild);
  }
}
