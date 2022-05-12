import { HttpClient } from './http';
import { MusicSearchResult } from './model/musicSearchResult';
import { MusicPlayerState } from './model/playerState';

const http = new HttpClient();

export async function initPlayerState(guildId: string): Promise<MusicPlayerState> {
  return http.get<MusicPlayerState>(`/guilds/${guildId}/player`);
}

export async function search(q: string): Promise<MusicSearchResult> {
  return http.get<MusicSearchResult>(`/music/search?q=${q}`);
}

export async function addToPlaylist(guildId: string, id: string, type: string): Promise<MusicPlayerState> {
  return http.post(`/guilds/${guildId}/player/playlist/add`, { id, type });
}

export async function removeFromPlaylist(guildId: string, index: number) {
  return http.post(`/guilds/${guildId}/player/playlist/remove`, { index });
}

export async function moveInPlaylist(guildId: string, from: number, to: number) {
  return http.post(`/guilds/${guildId}/player/playlist/move`, { from, to });
}

export async function play(guildId: string): Promise<void> {
  return http.post<void, MusicPlayerState>(`/guilds/${guildId}/player/play`);
}

export async function pause(guildId: string): Promise<void> {
  return http.post<void, MusicPlayerState>(`/guilds/${guildId}/player/pause`);
}

export async function skip(guildId: string): Promise<void> {
  return http.post<void, MusicPlayerState>(`/guilds/${guildId}/player/skip`);
}

export async function setTime(guildId: string, millis: number): Promise<void> {
  return http.post<void, {millis: number}>(`/guilds/${guildId}/player/time`, {millis});
}