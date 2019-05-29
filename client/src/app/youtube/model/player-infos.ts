import {PlayerStatus} from './player-status';
import {TrackInfos} from './track-infos';

export interface PlayerInfos {
  status: PlayerStatus;
  track?: TrackInfos;
}
