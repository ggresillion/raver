import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth/auth.service';
import {User} from '../models/user';
import {SoundService} from '../sound/sound.service';
import {GuildsService} from '../guilds/guilds.service';
import {Guild} from '../models/guild';

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

  constructor(
    private readonly authService: AuthService,
    private readonly soundService: SoundService,
    private readonly guildService: GuildsService,
  ) {
  }

  ngOnInit() {
    // this.authService.getConnectedUser().subscribe(user => this.connectedUser = user);
    this.guildService.getAvailableGuilds().subscribe(guilds => {
      this.guilds = guilds;
      if (guilds.length > 0) {
        this.selectedGuild = guilds[0];
      }
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
    this.selectedGuild = guild;
  }
}
