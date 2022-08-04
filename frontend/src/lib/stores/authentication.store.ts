import { writable, type Writable } from 'svelte/store';
import type { User } from '../model/user';

export const user: Writable<User> = writable();
