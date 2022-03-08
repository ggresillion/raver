import { configureStore, Store } from '@reduxjs/toolkit';
import logger from 'redux-logger';
import guildsReducer from './slices/guilds';
import musicPlayerReducer from './slices/music';

const store = configureStore({
    reducer: {
        musicPlayer: musicPlayerReducer,
        guilds: guildsReducer
    },
    middleware: getDefaultMiddleware =>
        getDefaultMiddleware()
            .concat(logger),
});

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;

export default store;