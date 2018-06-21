import {Injectable} from '@angular/core';
import {HttpClient, HttpEventType, HttpRequest, HttpResponse} from '@angular/common/http';
import Category from '../../model/Category';
import {environment} from '../../../environments/environment';
import {Observable, Subject} from "rxjs/index";

@Injectable({
  providedIn: 'root'
})
export class SongService {
  constructor(private http: HttpClient) {
  }

  getSongs(): Observable<Category[]> {
    return this.http.get<Category[]>(environment.apiEndpoint + '/songs');
  }

  playSong(song: string, category: string) {
    return this.http.get(environment.apiEndpoint + '/songs/play?song=' + song + '&category=' + category, {
      responseType: 'text'
    });
  }

  changeSongCategory(song: string, category: string) {
    return this.http.put(environment.apiEndpoint + '/songs?song=' + song + '&category=' + category, null, {
      responseType: 'text'
    });
  }

  public uploadSong(category: string, files: Set<File>): { [key: string]: Observable<number> } {
    const status = {};

    files.forEach(file => {
      const formData: FormData = new FormData();
      formData.append('category', category);
      formData.append('songs', file, file.name);

      const req = new HttpRequest('POST', environment.apiEndpoint + '/songs', formData, {
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
}
