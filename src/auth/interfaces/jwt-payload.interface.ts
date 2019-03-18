export interface JwtPayload {
  user: {
    id: string,
    email: string,
    username: string,
  };
}
