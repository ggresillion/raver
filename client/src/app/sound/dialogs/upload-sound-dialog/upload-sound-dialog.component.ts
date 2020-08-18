import {Component, ElementRef, Inject, ViewChild} from '@angular/core';
import {SoundService} from '../../sound.service';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';

@Component({
  selector: 'app-upload-sound-dialog',
  templateUrl: './upload-sound-dialog.component.html',
  styleUrls: ['./upload-sound-dialog.component.scss']
})
export class UploadSoundDialogComponent {

  @ViewChild('sound') 
  public soundInput: ElementRef;
  @ViewChild('image') 
  public imageInput: ElementRef;
  public soundName = '';
  public uploadProgress;

  constructor(public dialogRef: MatDialogRef<UploadSoundDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any,
              private songService: SoundService) {
  }

  public upload() {
    const sound = this.soundInput.nativeElement.files[0];
    const image = this.imageInput.nativeElement.files[0];
    this.uploadProgress = this.songService.uploadSound(this.soundName, this.data.categoryId, sound, image);
    this.uploadProgress.subscribe(null, null, () => {
      this.dialogRef.close();
    });
  }

  public onSoundSelection(e) {
    const file = e.target.files[0];
    if (!!file) {
      this.soundName = file.name.split('.').slice(0, -1).join('.');
      if (this.soundName.length > 30) {
        this.soundName = this.soundName.slice(0, 30);
      }
    }
  }
}
