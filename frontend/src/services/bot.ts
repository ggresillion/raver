import { HttpClient } from './http';

const http = new HttpClient();

export async function joinVoiceChannel(guildId: string): Promise<void> {
  return http.post(`/guilds/${guildId}/join`);
}

export async function getLatencyInMillis(): Promise<number> {
  const res = await http.get<{latency: number}>('/bot/latency');
  return res.latency;
}