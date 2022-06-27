<script>
  import Playlist from './Playlist.svelte';
  import Search from './Search.svelte';
</script>

<div class="music-container">
  <Search/>
</div>
<div>
  <Playlist/>
</div>

<style lang="scss">
  .music-container {
    display: flex;
    flex-direction: column;
    padding: 16px;
    max-height: 100%;
    box-sizing: border-box;
    gap: 16px;
    color: white;

    .side-title {
      font-family: "Roboto-Thin", serif;
    }
  }
</style>
