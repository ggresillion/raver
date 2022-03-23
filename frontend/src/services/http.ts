export class HttpClient {

  private static readonly API = 'http://localhost:8080/api';

  public async get<T>(path: string): Promise<T> {
    const res = await fetch(HttpClient.API + path, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
      },
    });

    if (res.status === 401) {
      window.location.href = HttpClient.API + '/auth/login';
    }

    const resBody = await res.json();
    return resBody;
  }

  public async post<R, B>(path: string, body?: B): Promise<R> {
    const res = await fetch(HttpClient.API + path, {
      method: 'POST',
      body: JSON.stringify(body),
      headers: {
        Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
      },
    });

    if (res.status === 401) {
      window.location.href = HttpClient.API + '/auth/login';
    }

    const resBody = await res.json();
    return resBody;
  }
}