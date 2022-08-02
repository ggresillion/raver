import type { User } from '../model/user';
import { get } from './http';

export async function getUser(): Promise<User> {
  return get<User>('/auth/user');
}
