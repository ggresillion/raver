import {Component, OnInit} from '@angular/core';
import {AuthService} from '../auth/auth.service';
import {User} from '../models/User';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {

  private connectedUser: User;

  constructor(private readonly authService: AuthService) {
  }

  ngOnInit() {
    this.authService.getConnectedUser().subscribe(user => this.connectedUser = user);
  }

}
