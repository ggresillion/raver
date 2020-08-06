import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { User } from '../models/user';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor(private readonly httpClient: HttpClient,
    private readonly router: Router) {
  }

  public setAccessToken(accessToken: string) {
    this.getStorage().setItem('accessToken', accessToken);
  }

  public setRefreshToken(refreshToken: string) {
    this.getStorage().setItem('refreshToken', refreshToken);
  }

  public getAccessToken() {
    return this.getStorage().getItem('accessToken');
  }

  public isConnected() {
    return !!this.getAccessToken();
  }

  public getConnectedUser(): Observable<User> {
    return this.httpClient.get<User>(environment.api + '/users/me');
  }

  public login() {
    window.location.href = environment.api + '/auth/login';
  }

  public logout() {
    this.getStorage().clear();
    this.router.navigate(['/login']);
  }

  private getStorage() {
    return localStorage;
  }
}
