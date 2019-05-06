import {Component, OnInit} from '@angular/core';
import {BotService} from '../../../shared/services/bot.service';
import {BotInfos} from '../../../shared/model/bot-infos';

@Component({
  selector: 'app-bot-status',
  templateUrl: './bot-status.component.html',
  styleUrls: ['./bot-status.component.scss']
})
export class BotStatusComponent implements OnInit {

  public botInfos: BotInfos;

  constructor(private readonly botService: BotService) {
  }

  ngOnInit() {
    this.botService.getInfos().subscribe(infos => this.botInfos = infos);
  }

}
