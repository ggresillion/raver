import { writable } from 'svelte/store';
import type { Toaster } from '../model/toaster';

const TOASTER_TIMEOUT = 3000;

export const toasters = writable<(Toaster & { id: number })[]>([]);

export function createToaster(toaster: Toaster) {
    const rand = Math.random();
    toasters.update(t => [...t, { ...toaster, id: rand }]);
    setTimeout(() => {
        toasters.update(t => {
            t.splice(t.findIndex(t2 => t2.id === rand), 1);
            return t;
        });
    }, TOASTER_TIMEOUT);
}
