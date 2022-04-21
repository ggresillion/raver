import { configureStore } from '@reduxjs/toolkit';
import logger from 'redux-logger';
import guildReducer from './slices/guild';
import musicPlayerReducer from './slices/music';

const store = configureStore({
  reducer: {
    musicPlayer: musicPlayerReducer,
    guild: guildReducer,
  },
  middleware: getDefaultMiddleware =>
    getDefaultMiddleware()
      .concat(logger),
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;