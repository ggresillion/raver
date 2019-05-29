import {PlayerStatus} from './player-status';
import {TrackInfos} from '../dto/track-infos';

export interface PlayerInfos {
  status: PlayerStatus;
  track?: TrackInfos;
}
