import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { Guild } from './model/guild';
import { config } from '../config';


export const guildAPI = createApi({
  reducerPath: 'guildAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: config.apiUrl, prepareHeaders: headers => {
      headers.set('Authorization', 'Bearer ' + localStorage.getItem('accessToken'));
      return headers;
    },
  }),
  endpoints: (builder) => ({
    getGuilds: builder.query<Guild[], void>({
      query: () => 'guilds',
    }),
    joinGuild: builder.query<void, string>({
      query: (id: string) => `guilds/${id}/join`,
    }),
  }),
});

export const { useGetGuildsQuery, useJoinGuildQuery } = guildAPI;
