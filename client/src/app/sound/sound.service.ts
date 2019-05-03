import {Injectable} from '@angular/core';
import {HttpClient, HttpEventType, HttpRequest} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {Observable, Subject} from 'rxjs';
import {Sound} from '../models/Sound';
import {VideoInfos} from './model/video-infos';

@Injectable({
  providedIn: 'root'
})
export class SoundService {

  private songsSubject = new Subject<Sound[]>();

  constructor(private http: HttpClient) {
  }

  public getSounds(): Observable<Sound[]> {
    this.refreshSounds();
    return this.songsSubject.asObservable();
  }

  public refreshSounds() {
    this.http.get<Sound[]>(`${environment.api}/sounds`)
      .subscribe(sounds => this.songsSubject.next(sounds));
  }

  public playSound(id: number) {
    return this.http.post(`${environment.api}/sounds/${id}/play`, null);
  }

  public uploadSound(name: string, categoryId: number, file: File): Observable<number> {
    const formData: FormData = new FormData();
    formData.append('name', name);
    if (categoryId) {
      formData.append('categoryId', categoryId.toString());
    }
    formData.append('sound', file, file.name);

    const req = new HttpRequest('POST', `${environment.api}/sounds`, formData, {
      reportProgress: true
    });

    const progress = new Subject<number>();
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
    return this.http.get<VideoInfos>(`${environment.api}/youtube/search?url=${videoURL}`);
  }

  public uploadFromYoutube(url: string, categoryId: number, name: string) {
    return this.http.post(`${environment.api}/youtube/upload`,
      {url, name, categoryId});
  }

  public changeCategory(soundId: number, categoryId: number) {
    return this.http.put(`${environment.api}/sounds/${soundId}`, {categoryId});
  }

  public deleteSound(soundId: number) {
    return this.http.delete(`${environment.api}/sounds/${soundId}`);
  }
}
