import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { config } from '../config';


export const botAPI = createApi({
  reducerPath: 'botAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: config.apiUrl, prepareHeaders: headers => {
      headers.set('Authorization', 'Bearer ' + localStorage.getItem('accessToken'));
      return headers;
    },
  }),
  endpoints: (builder) => ({
    joinVoiceChannel: builder.mutation<void, string>({
      query: (guildId: string) => ({
        url: `guilds/${guildId}/join`,
        method: 'post',
      }),
    }),
    getLatencyInMillis: builder.query<number, void>({
      query: () => 'bot/latency',
      transformResponse(baseQueryReturnValue: {latency: number}): number {
        return baseQueryReturnValue.latency;
      },
    }),
  }),
});

export const { useJoinVoiceChannelMutation, useGetLatencyInMillisQuery } = botAPI;
