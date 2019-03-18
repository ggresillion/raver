import {Component, Inject, OnInit} from '@angular/core';
import {CategoryService} from "../../services/category/category.service";
import {MAT_DIALOG_DATA, MatDialogRef} from "@angular/material";
import {AddFromYoutubeDialogComponent} from "../add-from-youtube-dialog/add-from-youtube-dialog.component";

@Component({
  selector: 'app-create-category-dialog',
  templateUrl: './create-category-dialog.component.html',
  styleUrls: ['./create-category-dialog.component.css']
})
export class CreateCategoryDialogComponent implements OnInit {

  public categoryName: string = "";

  constructor(public dialogRef: MatDialogRef<CreateCategoryDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any,
              private categoryService: CategoryService) {
  }

  ngOnInit() {
  }

  createCategory() {
    this.categoryService.createCategory(this.categoryName).subscribe(() => {
      this.dialogRef.close();
    });
  }
}
