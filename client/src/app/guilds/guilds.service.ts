import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Guild } from '../models/guild';
import { environment } from '../../environments/environment';
import { Observable, ReplaySubject } from 'rxjs';
import { tap } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class GuildsService {

  private guildsSubject = new ReplaySubject<Guild[]>(1);
  private guildSubject = new ReplaySubject<Guild>(1);

  constructor(private readonly http: HttpClient) {
    this.fetchGuilds();
  }

  public getAvailableGuilds(): Observable<Guild[]> {
    return this.guildsSubject.asObservable();
  }

  public getSelectedGuild(): Observable<Guild> {
    return this.guildSubject.asObservable();
  }

  public setSelectedGuild(guild: Guild): void {
    localStorage.setItem('defaultGuildId', guild.id);
    this.guildSubject.next(guild);
  }

  private fetchGuilds(): void {
    this.http.get<Guild[]>(environment.api + '/guilds')
      .pipe(tap(guilds => {
        guilds.sort(g1 => -g1.isBotInGuild);
      }))
      .subscribe(guilds => {
        this.guildsSubject.next(guilds);
        this.setDefaultGuild(guilds);
      });
  }

  private setDefaultGuild(guilds: Guild[]): void {
    let defaultGuild;
    const defaultGuildId = localStorage.getItem('defaultGuildId');
    if (defaultGuildId) {
      defaultGuild = guilds.find(g => g.id === defaultGuildId);
      if (!defaultGuild) {
        defaultGuild = guilds[0];
      }
    } else {
      defaultGuild = guilds[0];
    }
    this.guildSubject.next(defaultGuild);

  }
}
