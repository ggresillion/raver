import {Component, OnInit} from '@angular/core';
import {AuthService} from '../../services/auth/auth.service';
import {ActivatedRoute, Router} from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.scss']
})
export class LoginComponent implements OnInit {

  constructor(
    private readonly authService: AuthService,
    private readonly route: ActivatedRoute,
    private readonly router: Router,
  ) {
  }

  ngOnInit() {
    if (this.authService.isConnected()) {
      this.router.navigate(['/']);
    }
    this.route.queryParams.subscribe((params) => {
      const accessToken = params.access_token;
      const refreshToken = params.refresh_token;
      if (accessToken && refreshToken) {
        this.authService.setAccessToken(accessToken);
        this.authService.setRefreshToken(refreshToken);
        this.router.navigate(['/']);
      }
    });
  }

  public login() {
    this.authService.login();
  }
}
