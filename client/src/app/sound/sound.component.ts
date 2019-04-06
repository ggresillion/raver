import {Component, OnInit} from '@angular/core';
import {Sound} from '../models/Sound';
import {SoundService} from './sound.service';
import {MatDialog} from '@angular/material';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';
import {AddFromYoutubeDialogComponent} from './dialogs/add-from-youtube-dialog/add-from-youtube-dialog.component';

@Component({
  selector: 'app-sound',
  templateUrl: './sound.component.html',
  styleUrls: ['./sound.component.scss']
})
export class SoundComponent implements OnInit {

  public sounds: Sound[] = [];
  public categories = [{id: null, name: 'All'}];

  constructor(
    private readonly soundService: SoundService,
    public readonly dialog: MatDialog,
  ) {
  }

  public ngOnInit() {
    this.getSounds();
  }

  public onSoundClick(id: number) {
    this.soundService.playSound(id).subscribe(() => {
    });
  }

  public uploadSound(categoryId: number) {
    this.dialog.open(UploadSoundDialogComponent, {data: {categoryId}})
      .afterClosed().subscribe(this.getSounds);
  }

  public uploadFromYoutube(categoryId: number) {
    this.dialog.open(AddFromYoutubeDialogComponent, {data: {categoryId}})
      .afterClosed().subscribe(this.getSounds);
  }

  private getSounds() {
    this.soundService.getSounds().subscribe(sounds => this.sounds = sounds);
  }
}
