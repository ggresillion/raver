import { Injectable } from '@angular/core';
import { HttpClient, HttpEventType, HttpRequest } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { BehaviorSubject, Observable, Subject } from 'rxjs';
import { Sound } from '../models/sound';
import { VideoInfos } from './model/video-infos';
import { GuildsService } from '../guilds/guilds.service';
import { Guild } from '../models/guild';
import { map, flatMap, mergeMap, first } from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class SoundService {

  private sounds: Sound[] = [];
  private soundsSubject = new BehaviorSubject<Sound[]>([]);
  private selectedGuild: Guild;

  constructor(private http: HttpClient,
    private guildService: GuildsService) {
    this.guildService.getSelectedGuild()
      .subscribe(guild => this.selectedGuild = guild);
  }

  public getSounds(): Observable<Sound[]> {
    this.refreshSounds();
    return this.soundsSubject.asObservable();
  }

  public refreshSounds() {
    this.http.get<Sound[]>(`${environment.api}/sounds?guildId=${this.selectedGuild.id}`)
      .subscribe(sounds => {
        this.sounds = sounds;
        this.soundsSubject.next(sounds);
      });
  }

  public playSound(id: number): Observable<void> {
    return this.http.post<void>(`${environment.api}/sounds/${id}/play`, this.selectedGuild.id);
  }

  public uploadSound(name: string, categoryId: number, file: File): Observable<number> {
    const progress = new Subject<number>();
    const formData: FormData = new FormData();
    formData.append('name', name);
    formData.append('guildId', this.selectedGuild.id.toString());
    if (categoryId) {
      formData.append('categoryId', categoryId.toString());
    }
    formData.append('sound', file, file.name);

    const req = new HttpRequest('POST', `${environment.api}/sounds`, formData, {
      reportProgress: true
    });

    this.http.request(req).subscribe(event => {
      if (event.type === HttpEventType.UploadProgress) {
        const percentDone = Math.round(100 * event.loaded / event.total);
        progress.next(percentDone);
      } else if (event.type === 3) {
        progress.complete();
      }
    });

    return progress.asObservable();
  }

  public searchYoutube(videoURL: string) {
    return this.http.get<VideoInfos>(`${environment.api}/youtube/infos?url=${videoURL}`);
  }

  public uploadFromYoutube(url: string, categoryId: number, name: string): Observable<Sound> {
    return this.guildService.getSelectedGuild().pipe(mergeMap(guild => {
      console.log(guild)
      return this.http.post<Sound>(`${environment.api}/youtube/upload`,
        { url, name, categoryId, guildId: guild.id });
    }));
  }

  public changeCategory(soundId: number, categoryId: number) {
    return this.http.put(`${environment.api}/sounds/${soundId}`, { categoryId });
  }

  public deleteSound(soundId: number) {
    return this.http.delete(`${environment.api}/sounds/${soundId}`);
  }

  public setSearchString(search: string) {
    this.soundsSubject.next(this.sounds.filter(song => song.name.toLowerCase().includes(search.toLowerCase())));
  }
}
