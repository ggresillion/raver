import {Injectable} from '@angular/core';
import {environment} from '../../../environments/environment';
import {User} from '../../models/User';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor() {
  }

  public setAccessToken(accessToken: string) {
    this.getStorage().setItem('accessToken', accessToken);
  }

  public setRefreshToken(refreshToken: string) {
    this.getStorage().setItem('refreshToken', refreshToken);
  }

  public setConnectedUser(user: User) {
    this.getStorage().setItem('connectedUser', JSON.stringify(user));
  }

  public getAccessToken() {
    return this.getStorage().getItem('accessToken');
  }

  public isConnected() {
    return !!this.getAccessToken();
  }

  public getConnectedUser(): User {
    return JSON.parse(this.getStorage().getItem('connectedUser'));
  }

  public login() {
    window.location.href = environment.api + '/auth/login';
  }

  private getStorage() {
    return localStorage;
  }
}
