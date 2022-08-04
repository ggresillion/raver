<script>
    import Header from './_components/Header.svelte';
    import Guilds from './_components/Guilds.svelte';
    import Music from './_components/Music.svelte';
    import { selectedGuildId } from '$lib/stores/guild.store';
    import { onMount } from 'svelte';
    import Player from './_components/Player.svelte';
  
    onMount(() => {
      selectedGuildId.set(localStorage.getItem('selectedGuildId'));
    });
  </script>
  
  <style lang="scss">
    .header {
      position: fixed;
      height: 60px;
      width: 100%;
    }
  
    .guilds {
      position: absolute;
      top: 60px;
      left: 0;
      height: calc(100% - 160px);
      overflow-y: auto;
      overflow-x: hidden;
      width: 66px;
    }
  
    .main {
      position: fixed;
      top: 60px;
      left: 66px;
      height: calc(100% - 160px);
      width: calc(100% - 66px);
      background-color: rgb(56, 56, 56);
    }
  </style>
  
  <div class="header">
    <Header/>
  </div>
  <div class="guilds">
    <Guilds/>
  </div>
  
  <div class="main -z-10">
    <Music/>
    <Player/>
  </div>
  