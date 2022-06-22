import React, { ReactElement, useState } from 'react';
import ReactImageFallback from 'react-image-fallback';
import fallbackImage from '../../../../assets/icons/music_note.svg';
import searchIcon from '../../../../assets/icons/search_white_24dp.svg';
import playIcon from '../../../../assets/icons/play.svg';
import './Search.scss';
import { useSelector } from 'react-redux';
import { RootState } from '../../../../store';
import { useAddToPlaylistMutation, useLazySearchQuery } from '../../../../api/musicApi';
import { Loader } from '../../../../components/Loader';

export function Search() {

  const selectedGuild = useSelector((state: RootState) => state.selectedGuild);

  const [addToPlaylist] = useAddToPlaylistMutation();

  const [search, { data: results, isLoading: searchLoading, isSuccess: searchSuccess }] = useLazySearchQuery();

  const [searchQuery, setSearchQuery] = useState<string>('');

  async function onSearch(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e && e.key !== 'Enter') return;
    search(searchQuery);
  }

  async function onAddToPlaylist(id: string, type: string) {
    if (!selectedGuild) return;
    await addToPlaylist({ guildId: selectedGuild, id, type });
  }

  function getSongs(): ReactElement {
    if (!results || results?.tracks?.length < 1) return <div></div>;
    return <div className="type">
      <h2>Songs</h2>
      <div className="line">
        {results?.tracks.map(t => (
          <div key={t.id} className="card">
            <div className="thumbnail">
              <ReactImageFallback src={t.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>
            <span className="title">{t.title}</span>
            <span className="artist">{t.artists.map(a => a.name).join(', ')}</span>
            <button type="button"
              className="play"
              title="Play"
              style={{ backgroundImage: `url(${playIcon})` }}
              onClick={() => onAddToPlaylist(t.id, 'TRACK')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getArtists(): ReactElement {
    if (!results || results?.artists.length < 1) return <div></div>;
    return <div className="type">
      <h2>Artists</h2>
      <div className="line">
        {results?.artists.map(a => (
          <div key={a.id} className="card">
            <div className="thumbnail">
              <ReactImageFallback src={a.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>
            <span className="title">{a.name}</span>
            <button type="button"
              className="play"
              title="Play"
              style={{ backgroundImage: `url(${playIcon})` }}
              onClick={() => onAddToPlaylist(a.id, 'ARTIST')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getAlbums(): ReactElement {
    if (!results || results?.albums.length < 1) return <div></div>;
    return <div className="type">
      <h2>Albums</h2>
      <div className="line">
        {results?.albums.map(a => (
          <div key={a.id} className="card">
            <div className="thumbnail">
              <ReactImageFallback src={a.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>
            <span className="title">{a.name}</span>
            <span className="artist">{a.artists[0].name}</span>
            <button type="button"
              className="play"
              title="Play"
              style={{ backgroundImage: `url(${playIcon})` }}
              onClick={() => onAddToPlaylist(a.id, 'ALBUM')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getPlaylists(): ReactElement {
    if (!results || results?.playlists.length < 1) return <div></div>;
    return <div className="type">
      <h2>Playlists</h2>
      <div className="line">
        {results?.playlists.map(p => (
          <div key={p.id} className="card">
            <div className="thumbnail">
              <ReactImageFallback src={p.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>
            <span className="title">{p.name}</span>
            <button type="button"
              className="play"
              title="Play"
              style={{ backgroundImage: `url(${playIcon})` }}
              onClick={() => onAddToPlaylist(p.id, 'PLAYLIST')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getResults(): React.ReactElement {
    if (searchLoading) return <Loader/>;
    if (searchSuccess) return <div className="results">
      {getSongs()}
      {getArtists()}
      {getAlbums()}
      {getPlaylists()}
    </div>;
    return <div></div>;
  }


  return <>
    <div className="search-box">
      <input className="search-input"
        placeholder="Search for music"
        onKeyDown={onSearch}
        value={searchQuery}
        onChange={e => setSearchQuery(e.target.value)}/>
      <img className="search-icon" src={searchIcon} onClick={() => search(searchQuery)}></img>
    </div>
    {getResults()}
  </>;
}
