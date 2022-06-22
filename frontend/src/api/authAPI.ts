import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import { User } from './model/user';
import { config } from '../config';


export const authAPI = createApi({
  reducerPath: 'authAPI',
  baseQuery: fetchBaseQuery({
    baseUrl: config.apiUrl, prepareHeaders: headers => {
      headers.set('Authorization', 'Bearer ' + localStorage.getItem('accessToken'));
      return headers;
    },
  }),
  endpoints: (builder) => ({
    getUser: builder.query<User, void>({
      query: () => 'auth/user',
    }),
  }),
});

export const { useGetUserQuery } = authAPI;
