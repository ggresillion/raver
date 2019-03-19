import {Injectable} from '@angular/core';
import {environment} from '../../../environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  constructor() {
  }

  public isConnected() {
    return !!this.getConnectedUser();
  }

  public getConnectedUser() {
    return this.getStorage().getItem('connectedUser');
  }

  public login() {
    window.location.href = environment.api + '/auth/login';
  }

  private getStorage() {
    return localStorage;
  }
}
