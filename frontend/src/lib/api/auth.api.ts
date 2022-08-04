import { get } from "./http";

export async function logout() {
    return get<void>('/auth/logout');
  }
  