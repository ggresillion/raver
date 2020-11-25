import {BotStatus} from './bot-status.enum';

export class BotStateDTO {
  public readonly guilds: Array<{ status: BotStatus, id: string }>;
  public readonly volume: number;
}
