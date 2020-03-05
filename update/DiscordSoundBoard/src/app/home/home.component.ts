import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth/auth.service';
import {User} from '../models/user';
import {SoundService} from '../sound/sound.service';
import {GuildsService} from '../guilds/guilds.service';
import {Guild} from '../models/guild';
import {Location} from '@angular/common';
import {Router} from '@angular/router';
import {MatDialog} from '@angular/material/dialog';
import {AddBotToGuildDialogComponent} from '../guilds/dialogs/add-bot-to-guild-dialog/add-bot-to-guild-dialog.component';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  public connectedUser: User = {id: BigInt(169143950203027456), username: 'Un Cacatoès', avatar: 'a49fd2d7b4c48f6d2d070d210d08f69e'};
  public search: string;
  public guilds: Guild[] = [];
  public selectedGuild: Guild;
  public activatedRoute: any;

  constructor(
    private readonly authService: AuthService,
    private readonly soundService: SoundService,
    private readonly guildService: GuildsService,
    private readonly location: Location,
    private readonly router: Router,
    private readonly dialog: MatDialog,
  ) {
  }

  ngOnInit() {
    // this.authService.getConnectedUser().subscribe(user => this.connectedUser = user);
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

  public selectGuild(guild: Guild): void {
    if (!guild.isBotInGuild) {
      this.dialog.open(AddBotToGuildDialogComponent);
      return;
    }
    this.guildService.setSelectedGuild(guild);
  }

  public isLinkSelected(route: string): boolean {
    return this.activatedRoute === route;
  }

  public navigateTo(route: string): void {
    this.router.navigate([route]);
  }
}
