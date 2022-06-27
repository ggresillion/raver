<script lang="ts">
  import { playerState } from '$lib/stores/player-state.store.js';
  import { PlayerStatus } from '$lib/model/player-status.js';
  import { selectedGuildId } from '../../lib/stores/guild.store.js';
  import { subscribeToPlayerState } from '../../lib/api/music.api';

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

  selectedGuildId.subscribe((val) => {
    if (!val) return;
    subscribeToPlayerState(val);
  });
</script>

{#if !$playerState}
  <div>NOT OK</div>
{:else}
  <div class="player">
    <div class="track">
      <img class="track-thumbnail" src={$playerState.playlist[0].thumbnail}/>
      <div class="track-info">
        <span class="track-name">{$playerState.playlist[0].title}</span>
        <span class="track-artist">{$playerState.playlist[0].artists.map(a => a.name)
        .join(', ')}</span>
      </div>
      <!--        <Soundwave play={playerState.status === PlayerStatus.PLAYING}></Soundwave>-->
    </div>

    <div class="controls">
      <div class="buttons">
        <button type="button" style={{ visibility: 'hidden' }}>
        </button>
        <button type="button"
                class="skip-next"
                disabled={$playerState.playlist.length <= 0}></button>
      </div>

      {#if $playerState.status === PlayerStatus.PLAYING}
        <div class="progress-bar">
          <!--          <span>{getCurrentTime()}</span>-->
          <!--          <ReactSlider value={playerState.progress}-->
          <!--                       onAfterChange={setProgress}-->
          <!--                       min={0}-->
          <!--                       max={playerState?.playlist[0].duration}/>-->
          <span>{getTotalTime()}</span>
        </div>
      {:else}
        <div class="progress-bar">
          <span>-:-</span>
          <!--          <ReactSlider disabled/>-->
          <span>-:-</span>
        </div>
      {/if}
    </div>

    <div class="extra">
      <button class="join-channel"
              title="Join voice channel"
              disabled={!$selectedGuildId}></button>
    </div>
  </div>
{/if}
