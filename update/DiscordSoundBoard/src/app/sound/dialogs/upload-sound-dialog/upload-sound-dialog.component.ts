import {Component, ElementRef, Inject, ViewChild} from '@angular/core';
import {SoundService} from '../../sound.service';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';

@Component({
  selector: 'app-upload-sound-dialog',
  templateUrl: './upload-sound-dialog.component.html',
  styleUrls: ['./upload-sound-dialog.component.scss']
})
export class UploadSoundDialogComponent {

  @ViewChild('file') fileInput: ElementRef;
  public name = '';
  public uploadProgress;

  constructor(public dialogRef: MatDialogRef<UploadSoundDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any,
              private songService: SoundService) {
  }

  public upload() {
    const file = this.fileInput.nativeElement.files[0];
    this.uploadProgress = this.songService.uploadSound(this.name, this.data.categoryId, file);
    this.uploadProgress.subscribe(null, null, () => {
      this.dialogRef.close();
    });
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
