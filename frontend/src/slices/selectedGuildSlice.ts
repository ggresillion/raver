import { createSlice } from '@reduxjs/toolkit';

export const selectedGuildSlice = createSlice({
  name: 'selectedGuild',
  initialState: localStorage.getItem('selectedGuild'),
  reducers: {
    setSelectedGuild: (state, action) => {
      localStorage.setItem('selectedGuild', action.payload);
      return action.payload;
    },
  },
});

export const { setSelectedGuild } = selectedGuildSlice.actions;
