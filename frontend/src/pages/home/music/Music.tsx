import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import searchIcon from '../../../assets/icons/search_white_24dp.svg';
import playIcon from '../../../assets/icons/play_white_24dp.svg';
import addIcon from '../../../assets/icons/add_white_24dp.svg';
import { Loader } from '../../../components/Loader';
import { MusicSearchResult } from '../../../services/model/musicSearchResult';
import { Track } from '../../../services/model/track';
import { addToPlaylist, initPlayerState, search } from '../../../services/musicPlayer';
import { updatePlayerState } from '../../../slices/music';
import { RootState } from '../../../store';
import './Music.scss';

function secondsToTime(secs: number) {
  const hours = Math.floor(secs / (60 * 60));

  const divisor_for_minutes = secs % (60 * 60);
  const minutes = Math.floor(divisor_for_minutes / 60);

  const divisor_for_seconds = divisor_for_minutes % 60;
  const seconds = Math.ceil(divisor_for_seconds);

  return {
    hours,
    minutes,
    seconds: seconds < 10 ? '0' + seconds.toString() : seconds,
  };
}

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

  if (!ready) return (<Loader />);
  if (!selectedGuild) return (<div>Please select a Discord guild on the left pannel</div>);

  return (
    <div className='music-container'>
      <div className='search-box'>
        <input className='search-input' placeholder='Search for music' onKeyDown={onSearch} />
        <img className='search-icon' src={searchIcon}></img>
      </div>
      {results ? <div className='results'>
        <div className='type'>
          <h2>Songs</h2>
          <div className='line'>
            {results?.tracks.map(t => (
              <div key={t.id} className="card">
                <img src={t.thumbnail}></img>
                <span className='title'>{t.title}</span>
                <span className='artist'>{t.artists.map(a => a.name).join(', ')}</span>
                <button type='button'><img title='Add to playlist' src={addIcon} onClick={() => onAddToPlaylist(t.id, 'TRACK')}></img></button>
              </div>
            ))}
          </div>
        </div>
        <div className='type'>
          <h2>Artists</h2>
          <div className='line'>
            {results?.artists.map(a => (
              <div key={a.id} className="card">
                <img src={a.thumbnail}></img>
                <span className='title'>{a.name}</span>
                <button type='button'><img title='Play' src={playIcon} onClick={() => onAddToPlaylist(a.id, 'ARTIST')}></img></button>
              </div>
            ))}
          </div>
        </div>
        <div className='type'>
          <h2>Albums</h2>
          <div className='line'>
            {results?.albums.map(a => (
              <div key={a.id} className="card">
                <img src={a.thumbnail}></img>
                <span className='title'>{a.name}</span>
                <span className='artist'>{a.artists[0].name}</span>
                <button type='button'><img title='Play' src={playIcon} onClick={() => onAddToPlaylist(a.id, 'ALBUM')}></img></button>
              </div>
            ))}
          </div>
        </div>
        <div className='type'>
          <h2>Playlists</h2>
          <div className='line'>
            {results?.playlists.map(p => (
              <div key={p.id} className="card">
                <img src={p.thumbnail}></img>
                <span className='title'>{p.name}</span>
                <button type='button'><img title='Play' src={playIcon} onClick={() => onAddToPlaylist(p.id, 'PLAYLIST')}></img></button>
              </div>
            ))}
          </div>
        </div>
      </div> : <div></div>}
    </div>
  );
}