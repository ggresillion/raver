import React from 'react';
import { useSelector } from 'react-redux';
import { RootState } from '../../../store';
import { Playlist } from './components/Playlist';
import { Search } from './components/Search';
import './Music.scss';

export function Music() {

  const { selectedGuild } = useSelector((state: RootState) => state);

  if (!selectedGuild) return (<div>Please select a Discord guild on the left pannel</div>);

  return (
    <>
      <div className="music-container">
        <Search></Search>
      </div>
      <div>
        <Playlist></Playlist>
      </div>
    </>
  );
}
