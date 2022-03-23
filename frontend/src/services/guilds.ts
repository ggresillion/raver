import { Guild } from './model/guild';
import { webSocketClient } from './websocket';

export async function joinGuild(guild: Guild) {
  await webSocketClient.join(guild.id);
}