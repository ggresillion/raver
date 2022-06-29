import { get, post } from './http.api';
import type { MusicSearchResult } from '../model/music-search-result';
import { config } from './config.api';
import { playerState } from '../stores/player-state.store';
import type { MusicPlayerState } from '../model/player-state';
import { get as getValue } from 'svelte/store';

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

export async function moveInPlaylist(params: { guildId: string, from: number, to: number }) {
  return post<MusicPlayerState, { from: number, to: number }>(`/guilds/${params.guildId}/player/playlist/move`, {
    from: params.from,
    to: params.to,
  });
}

export async function removeFromPlaylist(params: { guildId: string, index: number }) {
  return post<MusicPlayerState, { index: number }>(`/guilds/${params.guildId}/player/playlist/remove`, {
    index: params.index,
  });
}

export async function subscribeToPlayerState(guildId: string) {
  const eventSource = new EventSource(`${config.apiUrl}/guilds/${guildId}/player/subscribe?access_token=${localStorage.getItem('accessToken')}`);
  eventSource.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('got message', data);
    playerState.set(data);
  }
}

export async function subscribeToProgress(guildId: string) {
  const eventSource = new EventSource(`${config.apiUrl}/guilds/${guildId}/player/progress/subscribe?access_token=${localStorage.getItem('accessToken')}`);
  eventSource.onmessage = (event) => {
    const data = event.data;
    const p = getValue(playerState);
    if (!p) return;
    p.progress = data;
    playerState.set(p);
  }
}
