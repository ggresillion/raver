import { writable } from 'svelte/store';
import type { MusicPlayerState } from '$lib/model/player-state';

export const playerState = writable<MusicPlayerState | null>();

