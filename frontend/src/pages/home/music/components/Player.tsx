import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import ReactSlider from 'react-slider';
import bufferingIcon from '../../../../assets/icons/loading.svg';
import next from '../../../../assets/icons/skip_next_white_24dp.svg';
import { joinVoiceChannel } from '../../../../services/bot';
import { PlayerStatus } from '../../../../services/model/playerStatus';
import { pause, play, setTime, skip } from '../../../../services/musicPlayer';
import { RootState } from '../../../../store';
import './Player.scss';
import { Soundwave } from './Soundwave';

function millisToMinutesAndSeconds(millis: number) {
  const minutes = Math.floor(millis / 60000);
  const seconds = ((millis % 60000) / 1000).toFixed(0);
  return {minutes, seconds};
}

export function Player() {

  const {playerState} = useSelector((state: RootState) => state.musicPlayer);
  const {selectedGuild} = useSelector((state: RootState) => state.guild);
  const [volume, setVolume] = useState<number>(80);

  function getButton() {
    if (!selectedGuild) return <button type='button' className='play' disabled></button>;
    switch (playerState.status) {
    case PlayerStatus.IDLE: 
    case PlayerStatus.PAUSED: 
      return <button type='button' className='play' onClick={() => play(selectedGuild?.id)}></button>;
    case PlayerStatus.PLAYING: 
      return <button type='button' className='pause' onClick={() => pause(selectedGuild?.id)}></button>;
    case PlayerStatus.BUFFERING: 
      return <button type='button' className='buffering' disabled></button>;
    default:
      return <button type='button' className='play' disabled></button>;
    }
  }

  function setProgress(progress: number) {
    if(!selectedGuild) return;
    setTime(selectedGuild.id, progress);
  }

  function getCurrentTime() {
    if (!playerState.progress) return '-:-';
    const time  = millisToMinutesAndSeconds(playerState.progress);
    const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
    return time.minutes + ':' + seconds;
  }

  function getTotalTime() {
    const time  = millisToMinutesAndSeconds(playerState?.playlist[0].duration);
    const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
    return time.minutes + ':' + seconds;
  }

  return (
    <div className="player">
      <div className="track">
        {playerState && playerState.playlist.length > 0 &&
                    <>
                      <img className='track-thumbnail' src={playerState.playlist[0].thumbnail} />
                      <div className='track-info'>
                        <span className='track-name'>{playerState.playlist[0].title}</span>
                        <span className='track-artist'>{playerState.playlist[0].artists.map(a => a.name).join(', ')}</span>
                      </div>
                      <Soundwave play={playerState.status === PlayerStatus.PLAYING}></Soundwave>
                    </>
        }
      </div>

      <div className="controls">
        <div className="buttons">
          <button type='button' style={{visibility: 'hidden'}}>
          </button>
          {getButton()}
          <button type='button' className='skip-next'  onClick={() => !selectedGuild || skip(selectedGuild?.id)} disabled={playerState.playlist.length <= 0}></button>
        </div>

        {playerState.status === PlayerStatus.PLAYING ? 
          <div className='progress-bar'>
            <span>{getCurrentTime()}</span>
            <ReactSlider value={playerState.progress} onAfterChange={setProgress} min={0} max={playerState?.playlist[0].duration}/>
            <span>{getTotalTime()}</span>
          </div> :
          <div className='progress-bar'>
            <span>-:-</span>
            <ReactSlider disabled/>
            <span>-:-</span>
          </div>
        }
      </div>

      <div className='extra'>
        <button className='join-channel' title='Join voice channel' disabled={!selectedGuild} onClick={() => !selectedGuild || joinVoiceChannel(selectedGuild?.id)}></button>
      </div>


      {/* <div className="volume">
        <img className="volume-down" src={volumeDown} />
        <ReactSlider className="progress-outer" thumbClassName='progress-button' defaultValue={volume} />
        <img className="volume-up" src={volumeUp} />
      </div> */}

    </div>
  );
}