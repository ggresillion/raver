import { Guild } from './model/guild';
import { webSocketClient } from './websocket';
import { HttpClient } from './http';

const http = new HttpClient();

export async function joinGuild(guild: Guild) {
  await webSocketClient.join(guild.id);
}

export async function getGuilds(): Promise<Guild[]> {
  const res = await http.get<{guilds: Guild[]}>('/bot/guilds');
  return res.guilds;
}

export function addGuild(): void {
  window.open('http://localhost:8080/api/bot/guilds/add');
}