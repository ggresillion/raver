import React from 'react';
import { useSelector } from 'react-redux';
import ReactSlider from 'react-slider';
import { PlayerStatus } from '../../../../api/model/playerStatus';
import { RootState } from '../../../../store';
import './Player.scss';
import { Soundwave } from './Soundwave';
import { useJoinVoiceChannelMutation } from '../../../../api/botAPI';
import {
  useGetMusicPlayerStateQuery,
  usePauseMutation,
  usePlayMutation,
  useSetTimeMutation,
  useSkipMutation,
} from '../../../../api/musicApi';

function millisToMinutesAndSeconds(millis: number) {
  const minutes = Math.floor(millis / 60000);
  const seconds = ((millis % 60000) / 1000).toFixed(0);
  return { minutes, seconds };
}

export function Player() {

  const selectedGuild = useSelector((state: RootState) => state.selectedGuild);
  const { data: playerState } = useGetMusicPlayerStateQuery(selectedGuild || '');

  const [pause] = usePauseMutation();
  const [play] = usePlayMutation();
  const [skip] = useSkipMutation();
  const [setTime] = useSetTimeMutation();
  const [joinChannel] = useJoinVoiceChannelMutation();

  function getButton() {
    if (!selectedGuild || !playerState) return <button type="button" className="play" disabled></button>;
    switch (playerState.status) {
    case PlayerStatus.IDLE:
    case PlayerStatus.PAUSED:
      return <button type="button" className="play" onClick={() => play(selectedGuild)}></button>;
    case PlayerStatus.PLAYING:
      return <button type="button" className="pause" onClick={() => pause(selectedGuild)}></button>;
    case PlayerStatus.BUFFERING:
      return <button type="button" className="buffering" disabled></button>;
    default:
      return <button type="button" className="play" disabled></button>;
    }
  }

  function setProgress(progress: number) {
    if (!selectedGuild) return;
    setTime({ guildId: selectedGuild, millis: progress });
  }

  function getCurrentTime() {
    if (!playerState || !playerState.progress) return '-:-';
    const time = millisToMinutesAndSeconds(playerState.progress);
    const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
    return time.minutes + ':' + seconds;
  }

  function getTotalTime() {
    if (!playerState) return;
    const time = millisToMinutesAndSeconds(playerState?.playlist[0].duration);
    const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
    return time.minutes + ':' + seconds;
  }

  if (!playerState) {
    return <div></div>;
  }

  return (
    <div className="player">
      <div className="track">
        {playerState && playerState.playlist.length > 0 &&
          <>
            <img className="track-thumbnail" src={playerState.playlist[0].thumbnail}/>
            <div className="track-info">
              <span className="track-name">{playerState.playlist[0].title}</span>
              <span className="track-artist">{playerState.playlist[0].artists.map(a => a.name).join(', ')}</span>
            </div>
            <Soundwave play={playerState.status === PlayerStatus.PLAYING}></Soundwave>
          </>
        }
      </div>

      <div className="controls">
        <div className="buttons">
          <button type="button" style={{ visibility: 'hidden' }}>
          </button>
          {getButton()}
          <button type="button"
            className="skip-next"
            onClick={() => !selectedGuild || skip(selectedGuild)}
            disabled={playerState.playlist.length <= 0}></button>
        </div>

        {playerState.status === PlayerStatus.PLAYING ?
          <div className="progress-bar">
            <span>{getCurrentTime()}</span>
            <ReactSlider value={playerState.progress}
              onAfterChange={setProgress}
              min={0}
              max={playerState?.playlist[0].duration}/>
            <span>{getTotalTime()}</span>
          </div> :
          <div className="progress-bar">
            <span>-:-</span>
            <ReactSlider disabled/>
            <span>-:-</span>
          </div>
        }
      </div>

      <div className="extra">
        <button className="join-channel"
          title="Join voice channel"
          disabled={!selectedGuild}
          onClick={() => !selectedGuild || joinChannel(selectedGuild)}></button>
      </div>


      {/* <div className="volume">
        <img className="volume-down" src={volumeDown} />
        <ReactSlider className="progress-outer" thumbClassName='progress-button' defaultValue={volume} />
        <img className="volume-up" src={volumeUp} />
      </div> */}

    </div>
  );
}
