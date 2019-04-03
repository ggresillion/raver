import {Component, OnInit} from '@angular/core';
import {Sound} from '../models/Sound';
import {SoundService} from './sound.service';
import {MatDialog} from '@angular/material';
import {UploadSoundDialogComponent} from './dialogs/upload-sound-dialog/upload-sound-dialog.component';

@Component({
  selector: 'app-sound',
  templateUrl: './sound.component.html',
  styleUrls: ['./sound.component.scss']
})
export class SoundComponent implements OnInit {

  public sounds: Sound[] = [];

  constructor(
    private readonly soundService: SoundService,
    public readonly dialog: MatDialog,
  ) {
  }

  ngOnInit() {
    this.soundService.getSounds().subscribe(sounds => this.sounds = sounds);
  }

  onSoundClick(id: number) {
    this.soundService.playSound(id).subscribe(() => {
    });
  }

  uploadSound() {
    this.dialog.open(UploadSoundDialogComponent);
  }
}
