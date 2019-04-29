import {Component, EventEmitter, OnInit, Output} from '@angular/core';
import {Sound} from '../models/Sound';
import {SoundService} from './sound.service';
import {MatDialog} from '@angular/material';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';
import {AddFromYoutubeDialogComponent} from './dialogs/add-from-youtube-dialog/add-from-youtube-dialog.component';
import {CreateCategoryDialogComponent} from './dialogs/create-category-dialog/create-category-dialog.component';
import {CategoryService} from '../category/category.service';
import {Category} from '../models/Category';
import {CdkDragDrop} from '@angular/cdk/typings/esm5/drag-drop';

@Component({
  selector: 'app-sound',
  templateUrl: './sound.component.html',
  styleUrls: ['./sound.component.scss']
})
export class SoundComponent implements OnInit {

  public sounds: Sound[] = [];
  public categories: Category[] = [];

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
      .afterClosed().subscribe(this.fetchSounds);
  }

  public uploadFromYoutube(categoryId: number) {
    this.dialog.open(AddFromYoutubeDialogComponent, {data: {categoryId}})
      .afterClosed().subscribe(this.fetchSounds);
  }

  public openCreateCategoryDialog() {
    this.dialog.open(CreateCategoryDialogComponent)
      .afterClosed().subscribe(this.fetchSounds);
  }

  public getSoundsForCategory(categoryId: number) {
    if (!categoryId) {
      return this.sounds.filter(sound => !sound.category);
    }
    return this.sounds.filter(sound => sound.category && sound.category.id === categoryId);
  }

  private fetchSounds() {
    this.categoryService.getCategories().subscribe(categories => {
      this.categories = categories;
      this.soundService.getSounds().subscribe(sounds => this.sounds = sounds);
    });
  }
}
