<script lang="ts">
  import deleteIcon from '$lib/assets/icons/delete_white_24dp.svg';
  import playlistIcon from '$lib/assets/icons/queue_music.svg';
  import arrowIcon from '$lib/assets/icons/double_arrow.svg';
  import { playerState } from '$lib/stores/player-state.store';
  import { dndzone } from 'svelte-dnd-action';
  import { flip } from 'svelte/animate';
  import { onMount } from 'svelte';
  import type { Track } from '$lib/model/track';
  import { moveInPlaylist, removeFromPlaylist } from '../../lib/api/music.api';
  import { selectedGuildId } from '../../lib/stores/guild.store';

  const flipDurationMs = 300;

  let isOpen = true;
  let savedUpcoming: Track[] & { key: number }[] = [];
  let upcoming: Track[] & { key: number }[] = [];

  onMount(() => {
    playerState.subscribe(s => {
      if (!s || s.playlist.length === upcoming.length + 1) return;
      upcoming = s.playlist.slice(1).map((t, i) => ({ ...t, key: i }));
      savedUpcoming = upcoming;
    });
  })

  function handleDndConsider(e) {
    upcoming = e.detail.items;
  }

  async function handleDndFinalize(e) {
    const id = e.detail.info.id;
    const from = savedUpcoming.findIndex(t => t.id === id) + 1;
    const to = e.detail.items.findIndex(t => t.id === id) + 1;
    try {
      await moveInPlaylist({ guildId: $selectedGuildId, from, to });
    } catch (e) {
      console.log('failed to move track', e);
      return;
    }
    upcoming = e.detail.items;
  }

  async function onDelete(id: string) {
    const index = upcoming.findIndex(t => t.id === id) + 1;
    await removeFromPlaylist({ index, guildId: $selectedGuildId });
  }

</script>

{#if $playerState}
  <div class="playlist" class:open={isOpen}>
    <div class="handle" on:click={() => isOpen = !isOpen}>
      <button class="icon open" style:background-image={`url(${arrowIcon})`}></button>
      <button class="icon closed" style:background-image={`url(${playlistIcon})`}></button>
    </div>
    <h2 class="title">Queue</h2>
    <h3 class="subtitle">Now playing</h3>
    <div class="track">
      <img class="track-thumbnail" src={$playerState.playlist[0].thumbnail} alt="thumbnail"/>
      <div class="track-info">
        <span class="track-name">{$playerState.playlist[0].title}</span>
        <span class="track-artist">{$playerState.playlist[0].artists.map(a => a.name)
        .join(', ')}</span>
      </div>
    </div>
    <h3 class="subtitle">Coming next</h3>
    <div class="tracks" use:dndzone={{items: upcoming, flipDurationMs}}
         on:consider={handleDndConsider}
         on:finalize={handleDndFinalize}>
      {#each upcoming as track(track.id)}
        <div class="track" animate:flip="{{duration: flipDurationMs}}">
          <img class="track-thumbnail" src={track.thumbnail} alt="thumbnail"/>
          <div class="track-info">
            <span class="track-name">{track.title}</span>
            <span class="track-artist">{track.artists.map(a => a.name).join(', ')}</span>
          </div>
          <button class="icon" style:background-image={`url(${deleteIcon})`}
                  on:click={() => onDelete(track.id)}></button>
        </div>
      {/each}
    </div>
  </div>
{/if}

<style lang="scss">
  .playlist {
    position: absolute;
    top: 0;
    right: -30vw;
    display: flex;
    flex-direction: column;
    height: 100%;
    width: 30vw;
    background: rgb(65, 65, 65);
    background: linear-gradient(90deg, rgba(65, 65, 65, 1) 0%, rgba(56, 56, 56, 1) 100%);
    padding: 6px;
    box-sizing: border-box;
    color: white;
    transition: all 0.2s ease-out;
    font-size: 14px;

    h2 {
      margin: 16px 0 0 0;
      font-size: 18px;
    }

    h3 {
      margin: 12px 0;
      font-size: 16px;
    }

    .icon {
      position: absolute;
      right: 2px;
      transition: all linear 0.2s;

      &.open {
        opacity: 0%;
        transform: rotate(180deg);
      }

      &.closed {
        opacity: 100%;
        transform: rotate(0);
      }
    }

    &.open {
      position: absolute;
      right: 0;

      .icon {
        position: absolute;

        &.open {
          opacity: 100%;
          transform: rotate(0);
        }

        &.closed {
          opacity: 0%;
          transform: rotate(180deg);
        }
      }
    }

    .handle {
      position: absolute;
      height: 36px;
      width: 34px;
      left: -34px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: rgb(65, 65, 65);
      border-top-left-radius: 12px;
      border-bottom-left-radius: 12px;
      top: calc(50% - 36px);
    }

    .title {
      font-family: "Roboto";
    }

    .subtitle {
      font-family: "Roboto-Thin";
    }

    .tracks {
      overflow-y: auto;
    }


  }

  .track {
    position: relative;
    display: flex;
    align-items: center;
    gap: 26px;
    color: white;
    padding: 6px;

    &:hover {
      background-color: rgba($color: #000000, $alpha: 0.2);
    }

    .track-thumbnail {
      height: 40px;
      filter: none;
    }

    .track-info {
      display: flex;
      flex-direction: column;
      justify-content: space-between;
      padding: 4px 0;

      .track-name {
        font-family: "Roboto-Thin";
      }

      .track-artist {
        font-family: "Roboto";
      }
    }

    &:hover {
      button {
        display: block;
      }
    }

    button {
      display: none;
      position: absolute;
      right: 24px;
      height: 24px;
      width: 24px;
      background-size: cover;
    }
  }
</style>
