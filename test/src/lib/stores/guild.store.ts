import { browser } from '$app/env';
import { writable } from 'svelte/store';

const STORE = 'selectedGuildId';

export const selectedGuildId = writable<string | null>();

if (browser) {
  selectedGuildId.set(localStorage.getItem(STORE));
}

selectedGuildId.subscribe(val => {
  if (browser) {
    if (!val) {
      localStorage.removeItem(STORE);
      return;
    }
    localStorage.setItem(STORE, val);
  }
});
