import { post } from './http.api';

export async function join(guildId: string) {
  return post<void, void>(`/guilds/${guildId}/join`);
}
