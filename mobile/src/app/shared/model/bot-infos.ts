import {BotStatus} from './bot-status';

export interface BotInfos {
  id: number;
  username: string;
  avatar: string;
  status: BotStatus;
}
