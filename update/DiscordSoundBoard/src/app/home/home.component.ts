import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth/auth.service';
import {User} from '../models/User';
import {SoundService} from '../sound/sound.service';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  public connectedUser: User;
  public search: string;

  constructor(
    private readonly authService: AuthService,
    private readonly soundService: SoundService,
  ) {
  }

  ngOnInit() {
    this.authService.getConnectedUser().subscribe(user => this.connectedUser = user);
  }

  public onSearchChange() {
    this.soundService.setSearchString(this.search);
  }

  public reset() {
    this.search = '';
    this.onSearchChange();
  }
}
