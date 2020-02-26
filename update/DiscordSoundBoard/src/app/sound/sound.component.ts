import {Component, OnInit} from '@angular/core';
import {Sound} from '../models/Sound';
import {SoundService} from './sound.service';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';
import {AddFromYoutubeDialogComponent} from './dialogs/add-from-youtube-dialog/add-from-youtube-dialog.component';
import {CreateCategoryDialogComponent} from './dialogs/create-category-dialog/create-category-dialog.component';
import {CategoryService} from '../shared/services/category.service';
import {Category} from '../models/Category';
import {animate, style, transition, trigger} from '@angular/animations';
import {RenameCategoryDialogComponent} from './dialogs/rename-category-dialog/rename-category-dialog.component';
import {MatDialog} from '@angular/material/dialog';

@Component({
  selector: 'app-sound',
  templateUrl: './sound.component.html',
  styleUrls: ['./sound.component.scss'],
  animations: [
    trigger('openClose', [
      transition(':enter', [
        style({transform: 'translateY(-100%)'}),
        animate('200ms ease-out')
      ]),
      transition(':leave', [
        style({transform: 'translateY(0)'}),
        animate('200ms ease-out',
          style({transform: 'translateY(-100%)'}))
      ]),
    ]),
  ],
})
export class SoundComponent implements OnInit {

  public sounds: Sound[] = [];
  public categories: Category[] = [];
  public isEditing = false;

  constructor(
    private readonly soundService: SoundService,
    private readonly categoryService: CategoryService,
    public readonly dialog: MatDialog,
  ) {
  }

  public ngOnInit() {
    this.fetchSounds();
  }

  public onSoundClick(id: number) {
    this.soundService.playSound(id).subscribe(() => {
    });
  }

  public uploadSound(categoryId: number) {
    this.dialog.open(UploadSoundDialogComponent, {data: {categoryId}})
      .afterClosed().subscribe(() => this.refreshSounds());
  }

  public uploadFromYoutube(categoryId: number) {
    this.dialog.open(AddFromYoutubeDialogComponent, {data: {categoryId}})
      .afterClosed().subscribe(() => this.refreshSounds());
  }

  public openCreateCategoryDialog() {
    this.dialog.open(CreateCategoryDialogComponent)
      .afterClosed().subscribe(() => this.fetchSounds());
  }

  public getSoundsForCategory(categoryId: number) {
    if (!categoryId) {
      return this.sounds.filter(sound => !sound.category);
    }
    return this.sounds.filter(sound => sound.category && sound.category.id === categoryId);
  }

  public deleteSound(soundId: number) {
    return this.soundService.deleteSound(soundId)
      .subscribe(() => this.refreshSounds());
  }

  public renameCategory(category: Category) {
    return this.dialog.open(RenameCategoryDialogComponent, {data: category})
      .afterClosed().subscribe(() => this.fetchSounds());
  }

  public deleteCategory(categoryId: number) {
    return this.categoryService.deleteCategory(categoryId)
      .subscribe(() => this.fetchSounds());
  }

  private fetchSounds() {
    this.categoryService.getCategories().subscribe(categories => {
      this.categories = categories;
      this.soundService.getSounds().subscribe(sounds => this.sounds = sounds);
    });
  }

  private refreshSounds() {
    this.soundService.refreshSounds();
  }
}
