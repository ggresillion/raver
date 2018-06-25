import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../../environments/environment";

@Injectable({
  providedIn: 'root'
})
export class CategoryService {

  constructor(private http: HttpClient) {
  }

  public createCategory(categoryName: string) {
    return this.http.post(environment.apiEndpoint + '/categories', {name: categoryName}, {
      responseType: 'text'
    });
  }
}
