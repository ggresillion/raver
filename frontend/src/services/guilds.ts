import { webSocketClient } from "./websocket";

export async function joinGuild(guildId: string) {
    await webSocketClient.join(guildId);
}