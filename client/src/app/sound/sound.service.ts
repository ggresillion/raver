import {Injectable} from '@angular/core';
import {HttpClient, HttpEventType, HttpRequest} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {Observable, Subject} from 'rxjs';
import {Sound} from '../models/Sound';

@Injectable({
  providedIn: 'root'
})
export class SoundService {
  constructor(private http: HttpClient) {
  }

  getSounds(): Observable<Sound[]> {
    return this.http.get<Sound[]>(environment.api + '/sounds');
  }

  playSound(id: number) {
    return this.http.post(`${environment.api}/sounds/${id}/play`, null);
  }

  changeSongCategory(song: string, category: string) {
    return this.http.put(environment.api + '/songs?song=' + song + '&category=' + category, null, {
      responseType: 'text'
    });
  }

  public uploadSong(category: string, files: Set<File>): { [key: string]: Observable<number> } {
    const status = {};

    files.forEach(file => {
      const formData: FormData = new FormData();
      formData.append('category', category);
      formData.append('songs', file, file.name);

      const req = new HttpRequest('POST', environment.api + '/songs', formData, {
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

      status[file.name] = {
        progress: progress.asObservable()
      };
    });

    return status;
  }

  public searchYoutube(videoURL: string) {
    return this.http.get(environment.api + '/youtube/search?url=' + videoURL);
  }

  public uploadFromYoutube(videoURL: string, category: string, name: string) {
    return this.http.get(environment.api + '/youtube/upload?url=' + videoURL
      + '&category=' + category
      + '&name=' + name, {
      responseType: 'text'
    });
  }
}
