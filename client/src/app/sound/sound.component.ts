import {Component, OnInit} from '@angular/core';
import {Sound} from '../models/Sound';
import {SoundService} from './sound.service';

@Component({
  selector: 'app-sound',
  templateUrl: './sound.component.html',
  styleUrls: ['./sound.component.scss']
})
export class SoundComponent implements OnInit {

  public sounds: Sound[] = [];

  constructor(private readonly songService: SoundService) {
  }

  ngOnInit() {
    this.songService.getSounds().subscribe(sounds => this.sounds = sounds);
  }

}
