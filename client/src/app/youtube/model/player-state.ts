import {PlayerStatus} from './player-status';
import {TrackInfos} from './track-infos';

export class PlayerState {
  status: PlayerStatus;
  playlist: TrackInfos[];
  totalLengthSeconds: number;
}
