// type MusicSearchResult struct {
// 	Tracks    []Track    `json:"tracks"`
// 	Playlists []Playlist `json:"playlists"`
// 	Albums    []Album    `json:"albums"`
// 	Artists   []Artist   `json:"artists"`
// }

import { Album } from './album';
import { Artist } from './artist';
import { Playlist } from './playlist';
import { Track } from './track';

export interface MusicSearchResult {
  tracks: Track[];
  playlists: Playlist[];
  albums: Album[];
  artists: Artist[];
}