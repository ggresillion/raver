import React, { ReactElement, useEffect, useState } from 'react';
import ReactImageFallback from 'react-image-fallback';
import { useDispatch, useSelector } from 'react-redux';
import addIcon from '../../../assets/icons/add_white_24dp.svg';
import fallbackImage from '../../../assets/icons/music_note.svg';
import playIcon from '../../../assets/icons/play_white_24dp.svg';
import searchIcon from '../../../assets/icons/search_white_24dp.svg';
import { Loader } from '../../../components/Loader';
import { MusicSearchResult } from '../../../services/model/musicSearchResult';
import { addToPlaylist, initPlayerState, search } from '../../../services/musicPlayer';
import { updatePlayerState } from '../../../slices/music';
import { RootState } from '../../../store';
import './Music.scss';

export function Music() {

  const dispatch = useDispatch();

  const [results, setResults] = useState<MusicSearchResult>();

  const {selectedGuild, ready} = useSelector((state: RootState) => state.guild);

  useEffect(() => {
    if (!selectedGuild) return;
    fetchPlayerState(selectedGuild.id);
  }, [selectedGuild]);

  async function fetchPlayerState(guildId: string) {
    dispatch(updatePlayerState(await initPlayerState(guildId)));
  }

  async function onSearch(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key !== 'Enter') return;
    const target = e.target as HTMLInputElement;
    const q = target.value;

    const result = await search(q);
    setResults(result);
  }

  async function onAddToPlaylist(id: string, type: string) {
    if (!selectedGuild) return;
    await addToPlaylist(selectedGuild.id, id, type);
  }

  function getSongs(): ReactElement {
    if (!results || results?.tracks.length < 1) return <div></div>;
    return <div className='type'>
      <h2>Songs</h2>
      <div className='line'>
        {results?.tracks.map(t => (
          <div key={t.id} className="card">
            <div className='thumbnail'>
              <ReactImageFallback src={t.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>
            <span className='title'>{t.title}</span>
            <span className='artist'>{t.artists.map(a => a.name).join(', ')}</span>
            <button type='button' className='play' title='Play' onClick={() => onAddToPlaylist(t.id, 'TRACK')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getArtists(): ReactElement {
    if (!results || results?.artists.length < 1) return <div></div>;
    return <div className='type'>
      <h2>Artists</h2>
      <div className='line'>
        {results?.artists.map(a => (
          <div key={a.id} className="card">
            <div className='thumbnail'>
              <ReactImageFallback src={a.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>           
            <span className='title'>{a.name}</span>
            <button type='button' className='play' title='Play' onClick={() => onAddToPlaylist(a.id, 'ARTIST')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getAlbums(): ReactElement {
    if (!results || results?.albums.length < 1) return <div></div>;
    return <div className='type'>
      <h2>Albums</h2>
      <div className='line'>
        {results?.albums.map(a => (
          <div key={a.id} className="card">
            <div className='thumbnail'>
              <ReactImageFallback src={a.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>    
            <span className='title'>{a.name}</span>
            <span className='artist'>{a.artists[0].name}</span>
            <button type='button' className='play' title='Play' onClick={() => onAddToPlaylist(a.id, 'ALBUM')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  function getPlaylists(): ReactElement {
    if (!results || results?.playlists.length < 1) return <div></div>;
    return <div className='type'>
      <h2>Playlists</h2>
      <div className='line'>
        {results?.playlists.map(p => (
          <div key={p.id} className="card">
            <div className='thumbnail'>
              <ReactImageFallback src={p.thumbnail} fallbackImage={fallbackImage}></ReactImageFallback>
            </div>    
            <span className='title'>{p.name}</span>
            <button type='button' className='play' title='Play' onClick={() => onAddToPlaylist(p.id, 'PLAYLIST')}></button>
          </div>
        ))}
      </div>
    </div>;
  }

  if (!ready) return (<Loader />);
  if (!selectedGuild) return (<div>Please select a Discord guild on the left pannel</div>);

  return (
    <div className='music-container'>
      <div className='search-box'>
        <input className='search-input' placeholder='Search for music' onKeyDown={onSearch} />
        <img className='search-icon' src={searchIcon}></img>
      </div>
      {results ? 
        <div className='results'>
          {getSongs()}
          {getArtists()}
          {getAlbums()}
          {getPlaylists()}
        </div> 
        : 
        <div></div>}
    </div>
  );
}