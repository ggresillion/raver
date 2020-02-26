import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-youtube-thumbnail',
  templateUrl: './youtube-thumbnail.component.html',
  styleUrls: ['./youtube-thumbnail.component.scss']
})
export class YoutubeThumbnailComponent implements OnInit {

  @Input()
  public id: string;

  constructor() {
  }

  ngOnInit(): void {
  }

}
