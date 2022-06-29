// type MusicSearchResult struct {
// 	Tracks    []Track    `json:"tracks"`
// 	Playlists []Playlist `json:"playlists"`
// 	Albums    []Album    `json:"albums"`
// 	Artists   []Artist   `json:"artists"`
// }

import type { Album } from './album';
import type { Artist } from './artist';
import type { Playlist } from './playlist';
import type { Track } from './track';

export interface MusicSearchResult {
  tracks: Track[];
  playlists: Playlist[];
  albums: Album[];
  artists: Artist[];
}
