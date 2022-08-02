import { config } from './config.api';

export async function get<T>(path: string) {
  return make<T, unknown>(path, 'GET');
}

export async function post<T, B>(path: string, body?: B) {
  return make<T, B>(path, 'POST', body);
}

async function make<T, B>(path: string, method: 'GET' | 'POST', body?: B) {
  const res = await fetch(config.apiUrl + path, {
    headers: {
      'Content-Type': 'application/json',
    },
    body: body ? JSON.stringify(body) : null,
    method,
  });
  if (!res.ok) {
    if (res.status === 401) {
      window.location.replace("/login");
    }
    throw res.body;
  }
  const obj = await res.json();
  return obj as T;
}
