import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Category } from '../../models/category';
import { Observable, ReplaySubject } from 'rxjs';
import { GuildsService } from '../../guilds/guilds.service';
import { flatMap, first, switchMap, tap } from 'rxjs/operators';

@Injectable()
export class CategoryService {

  public categoriesSubject = new ReplaySubject<Category[]>(1);

  constructor(
    private http: HttpClient,
    private guildsService: GuildsService) {
    this.fetchCategories();
  }

  public fetchCategories(): void {
    this.guildsService.getSelectedGuild().subscribe(guild => {
      return this.http.get<Category[]>(`${environment.api}/categories?guildId=${guild.id}`)
        .subscribe(categories => this.categoriesSubject.next(categories));
    });
  }

  public getCategories(): Observable<Category[]> {
    return this.categoriesSubject.asObservable();
  }

  public createCategory(categoryName: string) {
    return this.guildsService.getSelectedGuild()
      .pipe(switchMap(guild => {
        return this.http.post(`${environment.api}/categories`, { name: categoryName, guildId: guild.id }, {
          responseType: 'text'
        });
      }))
      .pipe(tap(() => {
        this.fetchCategories();
      }));
  }

  public renameCategory(categoryId: number, name: string) {
    return this.http.put(`${environment.api}/categories/${categoryId}`, { name })
      .pipe(tap(() => {
        this.fetchCategories();
      }));
  }

  public deleteCategory(categoryId: number) {
    return this.http.delete(`${environment.api}/categories/${categoryId}`)
      .pipe(tap(() => {
        this.fetchCategories();
      }));
  }
}
