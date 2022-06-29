import type { Album } from './album';
import type { Artist } from './artist';

export interface Track {
  id: string;
  title: string;
  artists: Artist[];
  album: Album;
  thumbnail: string;
  duration: number;
}
