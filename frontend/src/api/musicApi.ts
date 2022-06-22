import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { MusicSearchResult } from './model/musicSearchResult';
import { MusicPlayerState } from './model/playerState';
import { config } from '../config';

let eventSource: EventSource;

export const musicAPI = createApi({
  reducerPath: 'musicAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: config.apiUrl, prepareHeaders: headers => {
      headers.set('Authorization', 'Bearer ' + localStorage.getItem('accessToken'));
      return headers;
    },
  }),
  endpoints: (builder) => ({
    getMusicPlayerState: builder.query<MusicPlayerState, string>({
      query: (guildId: string) => `/guilds/${guildId}/player`,
      async onQueryStarted(guildId, { dispatch }) {
        if(eventSource) eventSource.close();
        eventSource = new EventSource(`${config.apiUrl}guilds/${guildId}/player/subscribe?access_token=${localStorage.getItem('accessToken')}`);
        eventSource.onmessage = (event) => {
          const data = event.data;
          dispatch(
            musicAPI.util.updateQueryData('getMusicPlayerState', guildId, (draft) => {
              Object.assign(draft, JSON.parse(data));
            }),
          );
        };
      },
    }),
    search: builder.query<MusicSearchResult, string>({
      query: (q: string) => `music/search?q=${q}`,
      keepUnusedDataFor: 0,
    }),
    addToPlaylist: builder.mutation<MusicPlayerState, { guildId: string, id: string, type: string }>({
      query: (params) => ({
        url: `/guilds/${params.guildId}/player/playlist/add`,
        method: 'post',
        body: { id: params.id, type: params.type },
      }),
    }),
    removeFromPlaylist: builder.mutation<MusicPlayerState, { guildId: string, index: number }>({
      query: (params) => ({
        url: `/guilds/${params.guildId}/player/playlist/remove`,
        method: 'post',
        body: { index: params.index },
      }),
    }),
    moveInPlaylist: builder.mutation<MusicPlayerState, { guildId: string, from: number, to: number }>({
      query: (params) => ({
        url: `/guilds/${params.guildId}/player/playlist/move`,
        method: 'post',
        body: { from: params.from, to: params.to },
      }),
    }),
    play: builder.mutation <void, string>({
      query: (guildId) => ({
        url: `/guilds/${guildId}/player/play`,
        method: 'post',
      }),
    }),
    pause: builder.mutation <void, string>({
      query: (guildId) => ({
        url: `/guilds/${guildId}/player/pause`,
        method: 'post',
      }),
    }),
    skip: builder.mutation <void, string>({
      query: (guildId) => ({
        url: `/guilds/${guildId}/player/skip`,
        method: 'post',
      }),
    }),
    setTime: builder.mutation <void, { guildId: string, millis: number }>({
      query: (params) => ({
        url: `/guilds/${params.guildId}/player/time`,
        method: 'post',
        body: { millis: params.millis },
      }),
    }),
  }),
});

export const {
  useGetMusicPlayerStateQuery,
  useLazySearchQuery,
  useAddToPlaylistMutation,
  useRemoveFromPlaylistMutation,
  useMoveInPlaylistMutation,
  usePlayMutation,
  usePauseMutation,
  useSkipMutation,
  useSetTimeMutation,
} = musicAPI;
