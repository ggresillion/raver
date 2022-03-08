import { MusicPlayerState } from "./model/playerState";
import { HttpClient } from "./http";

const http = new HttpClient();

export async function initPlayerState(guildId: string): Promise<MusicPlayerState> {
    return await http.get<MusicPlayerState>(`/guilds/${guildId}/player`);
}

export async function addToPlaylist(guildId: string, id: string): Promise<MusicPlayerState> {
    return await http.post(`/guilds/${guildId}/addToPlaylist`, { id });
}

export async function removeFromPlaylist(guildId: string, index: number) {
    return await http.post(`/guilds/${guildId}/removeFromPlaylist`, { index });
}

export async function moveInPlaylist(guildId: string, from: number, to: number) {
    return await http.post(`/guilds/${guildId}/moveInPlaylist`, { from, to });
}

export async function play(guildId: string): Promise<MusicPlayerState> {
    return await http.get<MusicPlayerState>(`/guilds/${guildId}/play`);
}

export async function pause(guildId: string): Promise<MusicPlayerState> {
    return await http.get<MusicPlayerState>(`/guilds/${guildId}/pause`);
}

export async function skip(guildId: string): Promise<MusicPlayerState> {
    return await http.get<MusicPlayerState>(`/guilds/${guildId}/skip`);
}