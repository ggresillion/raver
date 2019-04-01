import {Component, Inject, OnInit} from '@angular/core';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material';
import {SoundService} from '../../sound.service';

@Component({
  selector: 'app-upload-sound-dialog',
  templateUrl: './upload-sound-dialog.component.html',
  styleUrls: ['./upload-sound-dialog.component.scss']
})
export class UploadSoundDialogComponent implements OnInit {

  public name = '';

  constructor(public dialogRef: MatDialogRef<UploadSoundDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any,
              private songService: SoundService) {
  }

  ngOnInit() {
  }

  public upload() {
    this.dialogRef.close();
  }
}
