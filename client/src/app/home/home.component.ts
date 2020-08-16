import { Component, OnInit } from '@angular/core';
import { AuthService } from '../auth/auth.service';
import { User } from '../models/user';
import { SoundService } from '../sound/sound.service';
import { GuildsService } from '../guilds/guilds.service';
import { Guild } from '../models/guild';
import { Location } from '@angular/common';
import { Router } from '@angular/router';
import { BotService } from '../bot/bot.service';
import { environment } from '../../environments/environment';
import * as querystring from 'querystring';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  public connectedUser: User;
  public search: string;
  public guilds: Guild[];
  public selectedGuild: Guild;
  public activatedRoute: any;

  constructor(
    private readonly authService: AuthService,
    private readonly soundService: SoundService,
    private readonly guildService: GuildsService,
    private readonly location: Location,
    private readonly router: Router,
    private readonly botService: BotService,
  ) {
  }

  ngOnInit() {
    this.botService.getState().subscribe(s => console.log(s));
    this.authService.getConnectedUser().subscribe(user => this.connectedUser = user);
    this.guildService.getAvailableGuilds().subscribe(guilds => {
      this.guilds = guilds;
    });
    this.guildService.getSelectedGuild().subscribe(guild => {
      this.selectedGuild = guild;
    });
    this.activatedRoute = /[^/]*$/.exec(this.location.path())[0];
    this.location.onUrlChange(url => {
      this.activatedRoute = /[^/]*$/.exec(url)[0];
    });
  }

  public getGuilds(): Guild[] {
    return this.guilds.filter(g => g.isBotInGuild);
  }

  public onSearchChange() {
    this.soundService.setSearchString(this.search);
  }

  public reset() {
    this.search = '';
    this.onSearchChange();
  }

  public stringToRGB(str: string): string {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      // tslint:disable-next-line:no-bitwise
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    // tslint:disable-next-line:no-bitwise
    const c = (hash & 0x00FFFFFF)
      .toString(16)
      .toUpperCase();
    const color = '00000'.substring(0, 6 - c.length) + c;
    return '#' + color;
  }

  public selectGuild(guild: Guild) {
    this.guildService.setSelectedGuild(guild);
  }

  public addBotToGuild(): void {
    window.open(`${environment.discord.api}?${querystring.encode({
      client_id: environment.discord.clientId,
      permissions: environment.discord.permissions,
      redirect_uri: location.href,
      scope: environment.discord.scope,
    })}`);
  }

  public isLinkSelected(route: string): boolean {
    return this.activatedRoute === route;
  }

  public navigateTo(route: string): void {
    this.router.navigate([route]);
  }

  public disconnect() {
    this.authService.logout();
  }

  public joinMyChannel(): void {
    this.botService.joinMyChannel().subscribe(() => { });
  }
}
