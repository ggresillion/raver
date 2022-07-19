import { browser } from '$app/env';
import { writable } from 'svelte/store';

export const isAuthenticated = writable(browser && localStorage.getItem("accessToken"));
export const user = writable();
