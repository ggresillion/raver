<script lang="ts">
	import { getGuilds } from '$lib/api/bot.api';
	import { config } from '$lib/api/config.api';
	import addIcon from '$lib/assets/icons/add_white_24dp.svg';
	import Loader from '$lib/components/Loader.svelte';
	import { selectedGuildId } from '$lib/stores/guild.store';

	const guilds = getGuilds();

	function addGuild(): void {
		window.open(`${config.apiUrl}/guilds/add`);
	}

	function stringToRGB(str: string): string {
		let hash = 0;
		for (let i = 0; i < str.length; i++) {
			hash = str.charCodeAt(i) + ((hash << 5) - hash);
		}
		const c = (hash & 0x00ffffff).toString(16).toUpperCase();
		const color = '00000'.substring(0, 6 - c.length) + c;
		return '#' + color;
	}

	async function selectGuild(id: string) {
		selectedGuildId.update(() => id);
	}
</script>

<div
	class="bg-[#2a2a2a] flex flex-col align-center h-full w-full gap-2 p-2 overflow-y-auto box-border"
>
	{#await guilds}
		<Loader />
	{:then guilds}
		{#each guilds as { id, name, icon }, i}
			<div
				class="w-12 h-12 rounded-full cursor-pointer select-none flex items-center justify-center shadow-md before:absolute before:left-0 before:w-1 before:h-5 before:bg-white before:rounded-tr-full before:rounded-br-full"
				class:before:invisible={id !== $selectedGuildId}
				class:hover:scale-110={id !== $selectedGuildId}
				style:background={stringToRGB(name)}
				title={name}
				on:click={() => selectGuild(id)}
			>
				{#if icon}
					<img
						class="rounded-full"
						src={`https://cdn.discordapp.com/icons/${id}/${icon}`}
						alt={name}
					/>
				{:else}
					<span>{name.charAt(0).toUpperCase()}</span>
				{/if}
			</div>
		{/each}
		<div class="w-12 h-12 rounded-full cursor-pointer" title="Add to guild" on:click={addGuild}>
			<img class="w-12 h-12" src={addIcon} alt="Add to guild" />
		</div>
	{/await}
</div>
