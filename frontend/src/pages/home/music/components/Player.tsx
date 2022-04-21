import React, { useState } from 'react';
import { useSelector } from 'react-redux';
import pauseIcon from '../../../../assets/icons/pause_white_24dp.svg';
import playIcon from '../../../../assets/icons/play_white_24dp.svg';
import next from '../../../../assets/icons/skip_next_white_24dp.svg';
import previous from '../../../../assets/icons/skip_previous_white_24dp.svg';
import volumeDown from '../../../../assets/icons/volume_down_white_24dp.svg';
import volumeUp from '../../../../assets/icons/volume_up_white_24dp.svg';
import { joinVoiceChannel } from '../../../../services/bot';
import { PlayerStatus } from '../../../../services/model/playerStatus';
import { pause, play, skip } from '../../../../services/musicPlayer';
import { RootState } from '../../../../store';
import './Player.scss';
import { ProgressBar } from './ProgressBar';

export function Player() {

  const {playerState} = useSelector((state: RootState) => state.musicPlayer);
  const {selectedGuild} = useSelector((state: RootState) => state.guild);
  const [volume, setVolume] = useState<number>(80);

  function getButton() {
    if (!selectedGuild) return;
    switch (playerState.status) {
    case PlayerStatus.IDLE: 
    case PlayerStatus.PAUSED: 
      return <button type='button'><img className="play-pause" title='Play' src={playIcon} onClick={() => play(selectedGuild?.id)} ></img></button>;
    case PlayerStatus.PLAYING: 
      return <button type='button'><img className="play-pause" title='Pause' src={pauseIcon} onClick={() => pause(selectedGuild?.id)}></img></button>;
    default:
      return <button type='button' disabled><img className="play-pause" title='Play' src={playIcon}></img></button>;
    }
  }

  function setProgress(progress: number) {
    console.log(progress);
  }

  return (
    <div className="player">
      <div className="track">
        {playerState && playerState.playlist.length > 0 &&
                    <>
                      <img className='track-thumbnail' src={playerState.playlist[0].thumbnail} />
                      <div className='track-info'>
                        <span className='track-name'>{playerState.playlist[0].title}</span>
                        <span className='track-artist'>{playerState.playlist[0].artist}</span>
                      </div>
                    </>
        }
      </div>
      <div className="controls">

        <div className="buttons">
          <button type='button'>
            <img className="skip-previous" style={{opacity: 0}} title='Previous' src={previous}></img>
          </button>
          {getButton()}
          <button type='button' onClick={() => !selectedGuild || skip(selectedGuild?.id)}>
            <img className="skip-next" title='Next' src={next}></img>
          </button>
        </div>

        {playerState?.playlist.length > 0 ? 
          <div className='progress-bar'>
            <span>-:-</span>
            <ProgressBar progress={playerState?.progress || 0} onProgressChange={setProgress} disabled={true}/>
            <span>{playerState?.playlist[0].duration}</span>
          </div> :
          <div className='progress-bar'>
            <span>-:-</span>
            <ProgressBar disabled={true}/>
            <span>-:-</span>
          </div>
        }

      </div>

      {!selectedGuild || <button onClick={() => joinVoiceChannel(selectedGuild?.id)}>Join voice channel</button>}

      <div className="volume">
        <img className="volume-down" src={volumeDown} />
        <ProgressBar progress={volume} onProgressChange={setVolume} />
        <img className="volume-up" src={volumeUp} />
      </div>

    </div>
  );
}