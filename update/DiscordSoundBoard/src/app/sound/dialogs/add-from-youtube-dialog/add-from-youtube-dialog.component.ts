import {Component, Inject, OnInit} from '@angular/core';
import {SoundService} from '../../sound.service';
import {VideoInfos} from '../../model/video-infos';
import {MAT_DIALOG_DATA, MatDialogRef} from '@angular/material/dialog';

@Component({
  selector: 'app-add-from-youtube-dialog',
  templateUrl: './add-from-youtube-dialog.component.html',
  styleUrls: ['./add-from-youtube-dialog.component.scss']
})
export class AddFromYoutubeDialogComponent implements OnInit {

  public videoURL = '';
  public videoInfos: VideoInfos;
  public loading = false;
  public uploadName = '';

  constructor(public dialogRef: MatDialogRef<AddFromYoutubeDialogComponent>,
              @Inject(MAT_DIALOG_DATA) public data: any,
              private songService: SoundService) {
  }

  ngOnInit() {
  }

  public search() {
    this.loading = true;
    this.songService.searchYoutube(this.videoURL).subscribe((info) => {
      this.videoInfos = info;
      this.uploadName = info.title;
      this.loading = false;
    }, () => {
      this.loading = false;
    });
  }

  public upload() {
    this.loading = true;
    this.songService.uploadFromYoutube(this.videoURL, this.data.category,
      this.uploadName !== '' ? this.uploadName : this.videoInfos.title).subscribe(() => {
      this.dialogRef.close();
    });
  }
}
