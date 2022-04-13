import { createSlice, PayloadAction } from '@reduxjs/toolkit';
import { Guild } from '../services/model/guild';

interface GuildState {
  guilds?: Guild[];
  selectedGuild?: Guild;
  ready: boolean;
}

const initialState: GuildState = { 
  ready:false
};

export const guildsSlice = createSlice({
  name: 'guilds',
  initialState,
  reducers: {
    setGuilds: (state: GuildState, action: PayloadAction<{ guilds: Guild[] }>) => {
      return { ...state, guilds: action.payload.guilds };
    },
    setGuild: (state: GuildState, action: PayloadAction<{ selectedGuild?: Guild }>) => {
      return { ...state, ready: true, selectedGuild: action.payload.selectedGuild };
    },
  },
});

export const { setGuild, setGuilds } = guildsSlice.actions;

export default guildsSlice.reducer;