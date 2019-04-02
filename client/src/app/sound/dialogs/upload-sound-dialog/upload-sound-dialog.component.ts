import {Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import {MatDialogRef} from '@angular/material';
import {SoundService} from '../../sound.service';

@Component({
  selector: 'app-upload-sound-dialog',
  templateUrl: './upload-sound-dialog.component.html',
  styleUrls: ['./upload-sound-dialog.component.scss']
})
export class UploadSoundDialogComponent implements OnInit {

  public name = '';
  @ViewChild('file') fileInput: ElementRef;

  constructor(public dialogRef: MatDialogRef<UploadSoundDialogComponent>,
              private songService: SoundService) {
  }

  ngOnInit() {
  }

  public upload() {

    const file = this.fileInput.nativeElement.files[0];
    // this.songService.uploadSong(null, file);
    this.dialogRef.close();
  }

  public onFileSelection(e) {
    const file = e.target.files[0];
    if (!!file) {
      this.name = file.name.split('.').slice(0, -1).join('.');
      if (this.name.length > 30) {
        this.name = this.name.slice(0, 30);
      }
    }
  }
}
