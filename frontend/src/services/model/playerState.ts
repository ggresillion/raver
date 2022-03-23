import { PlayerStatus } from './playerStatus';
import { Track } from './track';

export interface MusicPlayerState {
  playlist: Track[];
  status: PlayerStatus;
}