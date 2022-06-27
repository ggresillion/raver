import { config } from './config.api';

export async function get<T>(path: string) {
  const res = await fetch(config.apiUrl + path, {
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('accessToken'),
    },
  });
  if (!res.ok) {
    throw res.body;
  }
  const obj = await res.json();
  return obj as T;
}

export async function post<T, B>(path: string, body: B) {
  const res = await fetch(config.apiUrl + path, {
    headers: {
      'Authorization': 'Bearer ' + localStorage.getItem('accessToken'),
      'Content-Type': 'application/json'
    },
    method: 'POST',
    body: JSON.stringify(body),
  });
  if (!res.ok) {
    throw res.body;
  }
  const obj = await res.json();
  return obj as T;
}
