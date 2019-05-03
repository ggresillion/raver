import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {SoundService} from '../../sound.service';
import {VideoInfos} from '../../model/video-infos';
import {CategoryService} from '../../../category/category.service';
import {Category} from '../../../models/Category';

@Component({
  selector: 'app-rename-category-dialog',
  templateUrl: './rename-category-dialog.component.html',
  styleUrls: ['./rename-category-dialog.component.scss']
})
export class RenameCategoryDialogComponent implements OnInit {

  public name = '';

  constructor(public dialogRef: MatDialogRef<RenameCategoryDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public category: Category,
              private categoryService: CategoryService) {
    this.name = category.name;
  }

  ngOnInit() {
  }

  public rename() {
    return this.categoryService.renameCategory(this.category.id, name)
      .subscribe(() => this.dialogRef.close());
  }
}
