<script lang="ts">
  import { selectedGuildId } from '$lib/stores/guild.store';
  import { getGuilds } from '$lib/api/guilds.api';
  import Loader from '$lib/components/Loader.svelte';
  import { config } from '$lib/api/config.api';
  import addIcon from '$lib/assets/icons/add_white_24dp.svg';

  const guilds = getGuilds();

  function addGuild(): void {
    window.open(`${config.apiUrl}guilds/add`);
  }

  function stringToRGB(str: string): string {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    const c = (hash & 0x00FFFFFF)
      .toString(16)
      .toUpperCase();
    const color = '00000'.substring(0, 6 - c.length) + c;
    return '#' + color;
  }

  async function selectGuild(id: string) {
    selectedGuildId.update(() => id);
  }

</script>

{#await guilds}
  <Loader/>
{:then guilds}
  <div class="guilds">
    {#each guilds as { id, name, icon }, i}
      <div class="guild" class:selected={id === $selectedGuildId}
           style:background={stringToRGB(name)}
           title={name}
           on:click={() => selectGuild(id)}>
        {#if icon}
          <img class="guild-icon"
               src={`https://cdn.discordapp.com/icons/${id}/${icon}`}
               alt={name}/>
        {:else}
          <span>{name.charAt(0).toUpperCase()}</span>
        {/if}
      </div>
    {/each}
    <div class="guild" title="Add to guild" onClick={() => addGuild()}>
      <img class="guild-icon" src={addIcon} alt="Add to guild"/>
    </div>
  </div>
{/await}

<style lang="scss">
  .guilds {
    background-color: #2a2a2a;
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 6px;
    gap: 6px;
    overflow-y: auto;
    box-sizing: border-box;

    .loader {
      transform: scale(0.5);
    }

    &::-webkit-scrollbar {
      display: none;
    }

    .guild {
      position: relative;
      width: 42px;
      min-width: 42px;
      height: 42px;
      min-height: 42px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      color: white;
      font-size: 24px;
      background-color: grey;
      cursor: pointer;
      user-select: none;
      box-shadow: 2px 2px 20px #0000002b;

      &.selected {
        &:before {
          content: "";
          position: absolute;
          left: -12px;
          height: 20px;
          width: 4px;
          background-color: white;
          border-top-right-radius: 5px;
          border-bottom-right-radius: 5px;
        }
      }

      &:hover {
        transform: scale(1.05);
      }

      .guild-icon {
        height: 42px;
        border-radius: 99px;
      }
    }
  }
</style>
