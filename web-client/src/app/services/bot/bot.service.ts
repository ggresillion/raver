import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../../environments/environment";
import {Observable} from "rxjs/index";

@Injectable({
  providedIn: 'root'
})
export class BotService {

  constructor(private http: HttpClient) {
  }

  getGuilds (): Observable<any> {
    return this.http.get(environment.apiEndpoint + '/manage/guilds');
  }
}
