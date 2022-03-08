import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { MusicPlayerState } from '../../services/model/playerState';

interface MusicState {
  playerState?: MusicPlayerState;
}

const initialState: MusicState = {};

export const musicPlayerSlice = createSlice({
  name: 'musicPlayer',
  initialState,
  reducers: {
    updatePlayerState: (_, action: PayloadAction<MusicPlayerState>) => {
      return { playerState: action.payload };
    }
  }
});

export const { updatePlayerState } = musicPlayerSlice.actions

export default musicPlayerSlice.reducer