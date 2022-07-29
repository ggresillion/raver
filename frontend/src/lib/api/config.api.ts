import { dev } from '$app/env';

export const config = {
    apiUrl: dev ? 'http://localhost:8080/api' : '/api',
};
