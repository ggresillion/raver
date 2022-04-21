import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { MusicPlayerState } from '../services/model/playerState';
import { PlayerStatus } from '../services/model/playerStatus';

interface MusicState {
  playerState: MusicPlayerState;
}

const initialState: MusicState = {
  playerState: {
    playlist: [],
    status: PlayerStatus.NOT_CONNECTED,
    progress: 0
  }
};

export const musicPlayerSlice = createSlice({
  name: 'musicPlayer',
  initialState,
  reducers: {
    updatePlayerState: (_, action: PayloadAction<MusicPlayerState>) => {
      return { playerState: action.payload };
    },
  },
});

export const { updatePlayerState } = musicPlayerSlice.actions;

export default musicPlayerSlice.reducer;