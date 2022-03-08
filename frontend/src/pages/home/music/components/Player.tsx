import { useState } from 'react';
import pause from '../../../../assets/icons/pause_white_24dp.svg';
import play from '../../../../assets/icons/play_white_24dp.svg';
import next from '../../../../assets/icons/skip_next_white_24dp.svg';
import previous from '../../../../assets/icons/skip_previous_white_24dp.svg';
import volumeDown from '../../../../assets/icons/volume_down_white_24dp.svg';
import volumeUp from '../../../../assets/icons/volume_up_white_24dp.svg';
import { MusicPlayerState } from '../../../../services/model/playerState';
import './Player.scss';
import { ProgressBar } from './ProgressBar';


enum PlayerStatus {
    PLAYING,
    PAUSED
}

export function Player() {

    const [progress, setProgress] = useState<number>(50);
    const [playerState] = useState<MusicPlayerState>();
    const [volume, setVolume] = useState<number>(80);

    return (
        <div className="player">
            <div className="track">
                {playerState &&
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
                    <img className="skip-previous" title='Previous' src={previous}></img>
                    {
                        playerState?.status === PlayerStatus.PAUSED ?
                            <img className="play-pause" title='Play' src={play} ></img> :
                            <img className="play-pause" title='Pause' src={pause} ></img>
                    }
                    <img className="skip-next" title='Next' src={next}></img>
                </div>

                <div className='progress-bar'>
                    <span>1:26</span>
                    <ProgressBar progress={progress} onProgressChange={setProgress} />
                    <span>3:01</span>
                </div>
            </div>

            <div className="volume">
                <img className="volume-down" src={volumeDown} />
                <ProgressBar progress={volume} onProgressChange={setVolume} />
                <img className="volume-up" src={volumeUp} />
            </div>

        </div>
    )
}