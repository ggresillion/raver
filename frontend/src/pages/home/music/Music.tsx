import React, { useEffect, useState } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import search from '../../../assets/icons/search_white_24dp.svg';
import { useAppSelector } from '../../../hooks';
import { Track } from '../../../services/model/track';
import { HttpClient } from '../../../services/http';
import { addToPlaylist, initPlayerState } from '../../../services/musicPlayer';
import { updatePlayerState } from '../../../slices/music';
import { RootState } from '../../../store';
import { Player } from './components/Player';
import { Playlist } from './components/Playlist';
import { Thumbnail } from './components/Thumbnail';
import './Music.scss';
import { Loader } from '../../../components/Loader';

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

  const [results, setResults] = useState<Track[]>();

  const {selectedGuild, ready} = useSelector((state: RootState) => state.guild);

  useEffect(() => {
    if (!selectedGuild) return;
    fetchPlayerState(selectedGuild.id);
  }, [selectedGuild]);

  async function fetchPlayerState(guildId: string) {
    dispatch(updatePlayerState(await initPlayerState(guildId)));
  }

  async function onEnter(e: React.KeyboardEvent<HTMLInputElement>) {
    if (e.key !== 'Enter') return;
    const target = e.target as HTMLInputElement;
    searchTracks(target.value);
  }

  async function searchTracks(q: string) {
    if (!q || q === '') return;
    const http = new HttpClient();
    const res = await http.get<Track[]>('/music/search?q=' + q);
    setResults(res);
  }

  async function addToPlaylistClick(guildId: string, track: Track) {
    await addToPlaylist(guildId, track.id);
  }

  if (!ready) return (<Loader />);
  if (!selectedGuild) return (<div>Please select a Discord guild on the left pannel</div>);

  return (
    <div className='music-container'>
      <div className='side-container'>
        <h2 className='side-title'>__ Search __</h2>
        <div className='search-box'>
          <input className='search-input' placeholder='Search for music' onKeyDown={onEnter} />
          <img className='search-icon' src={search}></img>
        </div>
        <div className='results'>
          {results?.map((r, i) => <div className='result' key={i} onClick={() => addToPlaylistClick(selectedGuild.id, r)}>
            <Thumbnail url={`https://img.youtube.com/vi/${r.id}/0.jpg`} />
            <div className='track-info'>
              <span className='track-name'>{r.title} - <span className='track-album'>{r.album}</span></span>
              <span className='track-artist'>{r.artist}</span>
            </div>
            <span className='track-duration'>{secondsToTime(r.duration).minutes}:{secondsToTime(r.duration).seconds}</span>
          </div>,
          )}
        </div>
      </div>
      <div className='side-container'>
        <Playlist />
      </div>
      <Player />
    </div>
  );
}