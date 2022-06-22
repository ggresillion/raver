import { configureStore } from '@reduxjs/toolkit';
import { guildAPI } from './api/guildAPI';
import { setupListeners } from '@reduxjs/toolkit/query';
import { botAPI } from './api/botAPI';
import { authAPI } from './api/authAPI';
import { musicAPI } from './api/musicApi';
import { selectedGuildSlice } from './slices/selectedGuildSlice';

const store = configureStore({
  reducer: {
    selectedGuild: selectedGuildSlice.reducer,
    [guildAPI.reducerPath]: guildAPI.reducer,
    [botAPI.reducerPath]: botAPI.reducer,
    [authAPI.reducerPath]: authAPI.reducer,
    [musicAPI.reducerPath]: musicAPI.reducer,
  },
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware()
      .concat(guildAPI.middleware)
      .concat(botAPI.middleware)
      .concat(authAPI.middleware)
      .concat(musicAPI.middleware),
});

setupListeners(store.dispatch);

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;
