import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Loader } from '../../../components/Loader';
import { initPlayerState } from '../../../services/musicPlayer';
import { updatePlayerState } from '../../../slices/music';
import { RootState } from '../../../store';
import { Playlist } from './components/Playlist';
import { Search } from './components/Search';
import './Music.scss';

export function Music() {

  const dispatch = useDispatch();

  const {selectedGuild, ready} = useSelector((state: RootState) => state.guild);

  useEffect(() => {
    if (!selectedGuild) return;
    fetchPlayerState(selectedGuild.id);
  }, [selectedGuild]);

  async function fetchPlayerState(guildId: string) {
    dispatch(updatePlayerState(await initPlayerState(guildId)));
  }

  if (!ready) return (<Loader />);
  if (!selectedGuild) return (<div>Please select a Discord guild on the left pannel</div>);

  return (
    <>
      <div className='music-container'>
        <Search></Search>
      </div>
      <div>
        <Playlist></Playlist>
      </div>
    </>
  );
}