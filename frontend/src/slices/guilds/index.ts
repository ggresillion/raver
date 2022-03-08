import { Action, createSlice, PayloadAction } from '@reduxjs/toolkit';

interface GuildState {
    loading: 'idle' | 'pending' | 'done' | 'failed';
    guildId?: string | null;
}

const initialState: GuildState = {
    loading: 'idle'
};

export const guildsSlice = createSlice({
    name: 'guilds',
    initialState,
    reducers: {
        setGuild: (_, action: PayloadAction<{ guildId: string | null }>) => {
            return { loading: 'done', guildId: action.payload.guildId };
        }
    }
});

export const { setGuild } = guildsSlice.actions;

export default guildsSlice.reducer