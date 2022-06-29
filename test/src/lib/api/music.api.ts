import { get, post } from './http.api';
import type { MusicSearchResult } from '../model/music-search-result';
import { config } from './config.api';
import { playerState } from '../stores/player-state.store';
import type { MusicPlayerState } from '../model/player-state';

export async function search(q: string) {
  return get<MusicSearchResult>(`/music/search?q=${q}`);
}

export async function addToPlaylist(params: { guildId: string, id: string, type: string }) {
  return post(`/guilds/${params.guildId}/player/playlist/add`, { id: params.id, type: params.type });
}

export async function getPlayerState(guildId: string) {
  return get<MusicPlayerState>(`/guilds/${guildId}/player`);
}

export async function play(guildId: string) {
  return post<MusicPlayerState, void>(`/guilds/${guildId}/player/play`);
}

export async function pause(guildId: string) {
  return post<MusicPlayerState, void>(`/guilds/${guildId}/player/pause`);
}

export async function skip(guildId: string) {
  return post<MusicPlayerState, void>(`/guilds/${guildId}/player/skip`);
}

export async function subscribeToPlayerState(guildId: string) {
  const eventSource = new EventSource(`${config.apiUrl}/guilds/${guildId}/player/subscribe?access_token=${localStorage.getItem('accessToken')}`);
  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('got message', data);
    playerState.set(data);
  }
}
