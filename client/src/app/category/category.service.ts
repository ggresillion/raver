import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {environment} from '../../environments/environment';
import {Category} from '../models/Category';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class CategoryService {

  constructor(private http: HttpClient) {
  }

  public getCategories(): Observable<Category[]> {
    return this.http.get<Category[]>(environment.api + '/categories');
  }

  public createCategory(categoryName: string) {
    return this.http.post(environment.api + '/categories', {name: categoryName}, {
      responseType: 'text'
    });
  }
}
