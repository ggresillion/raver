import { get } from './http';
import type { Guild } from '../model/guild';

export async function getGuilds(): Promise<Guild[]> {
  return get<Guild[]>('/guilds');
}
