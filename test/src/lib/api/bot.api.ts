import { get, post } from './http.api';

export async function getLatency() {
  return get<{latency: number}>(`/bot/latency`);
}

export async function join(guildId: string) {
  return post<void, void>(`/guilds/${guildId}/join`);
}
