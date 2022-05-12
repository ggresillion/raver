import { HttpClient } from './http';
import { User } from './model/user';

const http = new HttpClient();

export async function getUser(): Promise<User> {
  return http.get<User>('/auth/user');
}