<script lang="ts">
  import { playerState } from '$lib/stores/player-state.store';
  import { PlayerStatus } from '$lib/model/player-status';
  import { selectedGuildId } from '../../lib/stores/guild.store';
  import { getPlayerState, subscribeToPlayerState } from '../../lib/api/music.api';
  import RangeSlider from 'svelte-range-slider-pips';
  import { pause, play, skip } from '../../lib/api/music.api.js';
  import Loader from '../../lib/components/Loader.svelte';

  function millisToMinutesAndSeconds(millis: number) {
    const minutes = Math.floor(millis / 60000);
    const seconds = ((millis % 60000) / 1000).toFixed(0);
    return { minutes, seconds };
  }

  function getCurrentTime() {
    if (!$playerState || !$playerState.progress) return '-:-';
    const time = millisToMinutesAndSeconds($playerState.progress);
    const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
    return time.minutes + ':' + seconds;
  }

  function getTotalTime() {
    if (!playerState) return;
    const time = millisToMinutesAndSeconds($playerState?.playlist[0].duration);
    const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
    return time.minutes + ':' + seconds;
  }

  selectedGuildId.subscribe(async (val) => {
    if (!val) return;
    playerState.set(await getPlayerState($selectedGuildId));
    await subscribeToPlayerState(val);
  });
</script>

<div class="player">
  {#if !$playerState}
    <Loader/>
  {:else}
    <div class="track">
      {#if $playerState.playlist.length > 0}
        <img class="track-thumbnail" src={$playerState.playlist[0].thumbnail} alt="thumbnail"/>
        <div class="track-info">
          <span class="track-name">{$playerState.playlist[0].title}</span>
          <span class="track-artist">{$playerState.playlist[0].artists.map(a => a.name)
          .join(', ')}</span>
        </div>
        <!--        <Soundwave play={playerState.status === PlayerStatus.PLAYING}></Soundwave>-->
      {/if}
    </div>

    <div class="controls">
      <div>
        <button type="button" style="visibility: hidden"></button>
        {#if $playerState.status === PlayerStatus.IDLE || $playerState.status === PlayerStatus.PAUSED}
          <button type="button" class="play" on:click={play($selectedGuildId)}></button>
        {:else if $playerState.status === PlayerStatus.PLAYING}
          <button type="button" class="pause" on:click={pause($selectedGuildId)}></button>
        {:else if $playerState.status === PlayerStatus.BUFFERING}
          <button type="button" class="buffering" disabled></button>
        {:else}
          <button type="button" class="play" disabled></button>
        {/if}
        <button type="button"
                class="skip-next"
                on:click={skip($selectedGuildId)}
                disabled={!$playerState || $playerState.playlist.length <= 0}></button>
      </div>

      {#if $playerState && $playerState.status === PlayerStatus.PLAYING}
        <div>
          <!--          <span>{getCurrentTime()}</span>-->
          <!--          <ReactSlider value={playerState.progress}-->
          <!--                       onAfterChange={setProgress}-->
          <!--                       min={0}-->
          <!--                       max={playerState?.playlist[0].duration}/>-->
          <div class="progress-bar">
            <RangeSlider values={[50]}/>
          </div>
          <span>{getTotalTime()}</span>
        </div>
      {:else}
        <div>
          <div>-:-</div>
          <!--          <ReactSlider disabled/>-->
          <div class="progress-bar">
            <RangeSlider values={[0]} disabled/>
          </div>
          <span>-:-</span>
        </div>
      {/if}
    </div>

    <div class="extra">
      <button class="join-channel"
              title="Join voice channel"
              disabled></button>
    </div>
  {/if}
</div>

<style lang="scss">
  .player {
    position: fixed;
    bottom: 0;
    left: 0;
    height: 100px;
    width: 100vw;
    padding: 12px;
    background-color: rgb(53, 53, 53);
    filter: drop-shadow(3px 5px 2px rgb(0 0 0 / 0.4));
    color: white;
    box-sizing: border-box;

    display: grid;
    grid-template-columns: repeat(3, 1fr);
    grid-gap: 10px;

    .track {
      display: flex;
      align-items: center;
      gap: 26px;

      .track-thumbnail {
        height: 80px;
        filter: none;
      }

      .track-info {
        display: flex;
        flex-direction: column;
        justify-content: space-between;
        padding: 4px 0;

        .track-name {
          font-family: "Roboto-Thin", serif;
        }

        .track-artist {
          font-family: "Roboto", serif;
        }
      }
    }

    .controls {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 6px;

      & > div {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
        width: 100%;
      }

      .progress-bar {
        flex: 1;
      }
    }

    .extra {
      display: flex;
      justify-content: flex-end;
      align-items: center;

      .volume {
        width: 300px;
        display: flex;
        align-items: center;
        gap: 6px;

        .progress-inner {
          background: none;
          background-color: #03a9f4;
        }
      }

      .join-channel {
        height: 40px;
        width: 40px;
        border: none;
        background: url("$lib/assets/icons/join_channel.svg");
        background-size: cover;
        cursor: pointer;

        &:disabled {
          visibility: hidden;
        }

        &:hover:not(:disabled) {
          transform: scale(1.1);
        }
      }
    }
  }

  button {
    cursor: pointer;
    user-select: none;
    background-color: white;
    height: 36px;
    width: 36px;
    border-radius: 99px;
    border: none;
    display: flex;
    align-items: center;
    justify-content: center;
    background-size: contain;
    box-shadow: 0 0 20px 0 rgb(0 0 0 / 60%);
    -moz-box-shadow: 0 0 20px 0 rgb(0 0 0 / 60%);
    -webkit-box-shadow: 0 0 20px 0 rgb(0 0 0 / 60%);
    -o-box-shadow: 0 0 20px 0 rgb(0 0 0 / 60%);

    &.pause {
      background: url("$lib/assets/icons/pause.svg");
    }

    &.play {
      background: url("$lib/assets/icons/play.svg");
    }

    &.buffering {
      background: url("$lib/assets/icons/buffering.svg");
      animation: rotation 2s infinite linear;

      @keyframes rotation {
        from {
          transform: rotate(0deg);
        }
        to {
          transform: rotate(359deg);
        }
      }
    }

    &.skip-next {
      background: url("$lib/assets/icons/skip_next.svg");
      width: 32px;
      height: 32px;
    }

    &:disabled {
      opacity: 0.2;
      cursor: inherit;

      &:hover {
        transform: none;
      }
    }

    &:hover {
      transform: scale(1.1);
    }
  }

  //.slider {
  //  width: 100%;
  //  height: 10px;
  //  border: 1px rgb(27, 27, 27) solid;
  //  box-shadow: 0px 0px 25px black;
  //  border-radius: 99px;
  //  user-select: none;
  //  background: linear-gradient(90deg, #03a9f4 35%, rgba(109, 5, 255, 0.6194852941) 100%);
  //
  //  &:hover {
  //    .thumb {
  //      transform: scale(1.1);
  //    }
  //  }
  //
  //  .thumb {
  //    height: 12px;
  //    width: 12px;
  //    margin-top: -1px;
  //    background-color: white;
  //    border-radius: 99px;
  //    cursor: pointer;
  //  }
  //}
</style>
