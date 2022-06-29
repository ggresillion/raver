import type { PlayerStatus } from './player-status';
import type { Track } from './track';

export interface MusicPlayerState {
  playlist: Track[];
  status: PlayerStatus;
  progress?: number;
}
