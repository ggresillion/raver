import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-youtube-thumbnail',
  templateUrl: './youtube-thumbnail.component.html',
  styleUrls: ['./youtube-thumbnail.component.scss']
})
export class YoutubeThumbnailComponent implements OnInit {

  @Input()
  public link: string;

  constructor() {
  }

  ngOnInit(): void {
  }

  public getLink(): string {
    return this.link.replace('hqdefault', 'mqdefault');
  }

}
