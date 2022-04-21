import { HttpClient } from './http';

const http = new HttpClient();

export async function joinVoiceChannel(guildId: string): Promise<void> {
  return http.post(`/guilds/${guildId}/join`);
}