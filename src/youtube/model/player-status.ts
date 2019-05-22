import {Status} from './status';
import {TrackInfos} from '../dto/track-infos';

export interface PlayerStatus {
  status: Status;
  track: TrackInfos;
}
