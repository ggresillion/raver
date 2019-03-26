import {Component, OnInit} from '@angular/core';
import {AuthService} from '../../services/auth/auth.service';
import {User} from '../../../../../src/user/user.entity';

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
    // this.connectedUser = this.authService.getConnectedUser();
    console.log(this.connectedUser);
  }

}
