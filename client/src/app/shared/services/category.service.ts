import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Category } from '../../models/category';
import { Observable } from 'rxjs';
import { GuildsService } from '../../guilds/guilds.service';
import { Guild } from '../../models/guild';
import { flatMap, first } from 'rxjs/operators';

@Injectable()
export class CategoryService {

  constructor(
    private http: HttpClient,
    private guildsService: GuildsService
  ) {
  }

  public getCategories(): Observable<Category[]> {
    return this.guildsService.getSelectedGuild().pipe(flatMap(guild => {
      return this.http.get<Category[]>(`${environment.api}/categories?guildId=${guild.id}`);
    }));
  }

  public createCategory(categoryName: string) {
    return this.guildsService.getSelectedGuild()
      .pipe(flatMap(guild => {
        return this.http.post(`${environment.api}/categories`, { name: categoryName, guildId: guild.id }, {
          responseType: 'text'
        });
      }));
  }

  public renameCategory(categoryId: number, name: string) {
    return this.http.put(`${environment.api}/categories/${categoryId}`, { name });
  }

  public deleteCategory(categoryId: number) {
    return this.http.delete(`${environment.api}/categories/${categoryId}`);
  }
}
