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

  constructor(private http: HttpClient,
    private guildService: GuildsService) {
      this.refreshSounds();
  }

  public getSounds(): Observable<Sound[]> {
    return this.soundsSubject.asObservable();
  }

  public refreshSounds() {
    return this.guildService.getSelectedGuild().subscribe(guild => {
      return this.http.get<Sound[]>(`${environment.api}/sounds?guildId=${guild.id}`)
        .subscribe(sounds => {
          this.sounds = sounds;
          this.soundsSubject.next(sounds);
        });
    });
  }

  public playSound(id: number): Observable<void> {
    return this.guildService.getSelectedGuild().pipe(first()).pipe(flatMap(guild => {
      return this.http.post<void>(`${environment.api}/sounds/${id}/play`, guild.id);
    }));
  }

  public uploadSound(name: string, categoryId: number, sound: File, image?: File): Observable<number> {
    return this.guildService.getSelectedGuild().pipe(first()).pipe(flatMap(guild => {
      const progress = new Subject<number>();
      const formData: FormData = new FormData();
      formData.append('name', name);
      formData.append('guildId', guild.id.toString());
      if (categoryId) {
        formData.append('categoryId', categoryId.toString());
      }
      formData.append('sound', sound, sound.name);
      if (!!image) {
        formData.append('image', image);
      }

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
    }));
  }

  public searchYoutube(videoURL: string) {
    return this.http.get<VideoInfos>(`${environment.api}/youtube/infos?url=${videoURL}`);
  }

  public uploadFromYoutube(url: string, categoryId: number, name: string): Observable<Sound> {
    return this.guildService.getSelectedGuild().pipe(first()).pipe(mergeMap(guild => {
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
