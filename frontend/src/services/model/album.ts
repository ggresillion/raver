import { Artist } from './artist';

export interface Album {
  id: string;
  name: string;
  thumbnail: string;
  artists: Artist[];
}