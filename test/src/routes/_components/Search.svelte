<script lang="ts">
  import searchIcon from '$lib/assets/icons/search_white_24dp.svg';
  import playIcon from '$lib/assets/icons/play.svg';
  import { search, addToPlaylist } from '$lib/api/music.api';
  import type { MusicSearchResult } from '$lib/model/music-search-result';
  import Loader from '$lib/components/Loader.svelte';
  import { selectedGuildId } from '$lib/stores/guild.store';

  let searchInput = '';
  let searchPromise: Promise<MusicSearchResult>;

  function handleKeyUp(e: KeyboardEvent) {
    if (e.code === 'Enter') {
      searchPromise = search(searchInput);
    }
  }

  function onAddToPlaylist(id: string, type: string) {
    addToPlaylist({ guildId: $selectedGuildId, id, type });
  }
</script>

<div class="search-box">
  <input class="search-input"
         placeholder="Search for music"
         bind:value={searchInput}
         on:keyup={handleKeyUp}
  />
  <img class="search-icon" src={searchIcon} onClick={() =>  search()} alt="Search"/>
</div>

{#if !!searchPromise}
  {#await searchPromise}
    <Loader/>
  {:then results}
    <div class="results">

      <div class="type">
        <h2>Songs</h2>
        <div class="line">
          {#each results.tracks as t}
            <div class="card">
              <div class="thumbnail">
                <img src={t.thumbnail} alt="Thumbnail"/>
              </div>
              <span class="title">{t.title}</span>
              <span class="artist">{t.artists.map(a => a.name).join(', ')}</span>
              <button type="button"
                      class="play"
                      title="Play"
                      style:background-image={`url(${playIcon})`}
                      on:click={() => onAddToPlaylist(t.id, 'TRACK')}></button>
            </div>
          {/each}
        </div>
      </div>

      <div class="type">
        <h2>Artists</h2>
        <div class="line">
          {#each results.artists as t}
            <div class="card">
              <div class="thumbnail">
                <img src={t.thumbnail} alt="Thumbnail"/>
              </div>
              <span class="title">{t.name}</span>
              <button type="button"
                      class="play"
                      title="Play"
                      style:background-image={`url(${playIcon})`}
                      on:click={() => onAddToPlaylist(t.id, 'ARTIST')}></button>
            </div>
          {/each}
        </div>
      </div>

      <div class="type">
        <h2>Albums</h2>
        <div class="line">
          {#each results.albums as t}
            <div class="card">
              <div class="thumbnail">
                <img src={t.thumbnail} alt="Thumbnail"/>
              </div>
              <span class="title">{t.name}</span>
              <span class="artist">{t.artists.map(a => a.name).join(', ')}</span>
              <button type="button"
                      class="play"
                      title="Play"
                      style:background-image={`url(${playIcon})`}
                      on:click={() => onAddToPlaylist(t.id, 'ALBUM')}></button>
            </div>
          {/each}
        </div>
      </div>

      <div class="type">
        <h2>Playlists</h2>
        <div class="line">
          {#each results.playlists as t}
            <div class="card">
              <div class="thumbnail">
                <img src={t.thumbnail} alt="Thumbnail"/>
              </div>
              <span class="title">{t.name}</span>
              <button type="button"
                      class="play"
                      title="Play"
                      style:background-image={`url(${playIcon})`}
                      on:click={() => onAddToPlaylist(t.id, 'PLAYLIST')}></button>
            </div>
          {/each}
        </div>
      </div>

    </div>
  {:catch error}
    <span>{error}</span>
  {/await}
{/if}

<style lang="scss">
  .search-box {
    position: relative;
    height: 40px;
    box-sizing: border-box;
    width: 30vw;

    .search-input {
      padding: 0 12px;
      height: 37px;
      width: calc(100% - 35px);
      border-radius: 25px 0 0 25px;
      margin: 0;
      box-sizing: border-box;
      border: 2px solid #03a9f4;

      &:focus {
        outline: none;
      }
    }

    .search-icon {
      position: absolute;
      right: -1px;
      top: 0;
      height: 37px;
      border-top-right-radius: 25px;
      border-bottom-right-radius: 25px;
      background-color: #03a9f4;
      cursor: pointer;
      user-select: none;
    }
  }

  .results {
    display: flex;
    flex-direction: column;
    overflow-y: auto;

    .type {
      .line {
        display: flex;
        flex-direction: row;
        flex-flow: row wrap;
        gap: 6px;
        overflow-x: auto;
        overflow-y: hidden;
        padding: 6px 0;
        height: 190px;

        .card {
          position: relative;
          display: flex;
          flex-direction: column;
          background-color: #2c2c2c;
          border-radius: 6px;
          padding: 6px;
          gap: 6px;
          max-height: 180px;
          min-height: 180px;
          max-width: 120px;
          min-width: 120px;
          cursor: pointer;

          span {
            font-size: 14px;
            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
          }

          .title {
          }

          .artist {
            opacity: 0.6;
          }

          .thumbnail {
            width: 120px;
            height: 120px;
            background: linear-gradient(to top right, #03a9f4 35%, rgba(109, 5, 255, 0.6194852941) 100%);
            border-radius: 5px;

            img {
              width: 100%;
              height: 100%;
              object-fit: cover;
              border-radius: 5px;
            }
          }

          &:hover {
            button {
              opacity: unset;
            }
          }

          button {
            position: absolute;
            opacity: 0;
            bottom: 0;
            right: 0;
            cursor: pointer;
            user-select: none;
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
            background-color: transparent;

            &.play {
              margin-right: 6px;
              margin-bottom: 6px;
            }

            &:hover {
              transform: scale(1.1);
            }
          }
        }
      }
    }
  }
</style>
