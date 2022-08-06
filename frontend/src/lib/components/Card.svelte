<script lang="ts">
	import fallbackImage from '$lib/assets/icons/music_note.svg';
	import playIcon from '$lib/assets/icons/play.svg';
	import { selectedGuildId } from '$lib/stores/guild.store';
	import { addToPlaylist } from '../api/music.api';

	export let id: string;
	export let thumbnail: string;
	export let title: string;
	export let subtitle: string;
	export let type: 'TRACK' | 'ARTIST' | 'ALBUM' | 'PLAYLIST';

	function onAddToPlaylist(id: string) {
		if (!$selectedGuildId) return;
		addToPlaylist({ guildId: $selectedGuildId, id, type });
	}

	function handleImageError(ev: any) {
		return (ev.target.src = fallbackImage);
	}
</script>

<div
	class="relative flex flex-col rounded-lg cursor-pointer mb-2 h-full w-36 drop-shadow-2xl bg-[#383838] group"
>
	<div class="w-36 h-36">
		<img
			class="rounded-t-lg min-h-full max-h-full min-w-full max-w-full"
			src={thumbnail}
			on:error={handleImageError}
			alt="Thumbnail"
		/>
	</div>
	<div class="p-2 overflow-hidden">
		<span class="max-h-20 text-ellipsis whitespace-nowrap">{title}</span><br />
		<span class="opacity-60 max-h-20 text-ellipsis overflow-hidden whitespace-nowrap"
			>{subtitle}</span
		>
	</div>
	<button
		type="button"
		class="fixed right-2 bottom-2 h-10 w-10 bg-cover invisible group-hover:visible hover:scale-110"
		title="Play"
		style:background-image={`url(${playIcon})`}
		on:click={() => onAddToPlaylist(id)}
	/>
</div>
