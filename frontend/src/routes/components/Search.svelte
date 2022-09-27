<script lang="ts">
	import { search } from '$lib/api/music.api';
	import Loader from '$lib/components/Loader.svelte';
	import type { MusicSearchResult } from '$lib/model/music-search-result';
	import Card from '$lib/components/Card.svelte';
	import SearchBar from '$lib/components/SearchBar.svelte';

	let searchPromise: Promise<MusicSearchResult>;

	function submitSearch(e: any): void {
		e.preventDefault();
		searchPromise = search(e.detail);
	}

</script>

<SearchBar on:submit={submitSearch} />

{#if !!searchPromise}
	{#await searchPromise}
		<Loader />
	{:then results}
		<div class="flex flex-col overflow-y-auto mt-2 mb-2">
			<!-- Songs -->
			<div class="text-white m-2">
				<h2 class="text-xl mb-4">Songs</h2>
				<div class="flex flex-row overflow-x-auto overflow-y-hidden gap-2">
					{#each results.tracks as track}
						<Card
							id={track.id}
							title={track.title}
							subtitle={track.artists.map((a) => a.name).join(', ')}
							thumbnail={track.thumbnail}
							type="TRACK"
						/>
					{/each}
				</div>
			</div>

			<!-- Artists -->
			<div class="text-white m-2">
				<h2 class="text-xl">Artists</h2>
				<div class="flex flex-row overflow-x-auto overflow-y-hidden gap-2">
					{#each results.artists as artist}
						<Card
							id={artist.id}
							title={artist.name}
							subtitle={'Artist'}
							thumbnail={artist.thumbnail}
							type="ARTIST"
						/>
					{/each}
				</div>
			</div>

			<!-- Albums -->
			<div class="text-white m-2">
				<h2 class="text-xl">Albums</h2>
				<div class="flex flex-row overflow-x-auto overflow-y-hidden gap-2">
					{#each results.albums as album}
						<Card
							id={album.id}
							title={album.name}
							subtitle={album.artists.map((a) => a.name).join(', ')}
							thumbnail={album.thumbnail}
							type="ALBUM"
						/>
					{/each}
				</div>
			</div>

			<!-- Playlists -->
			<div class="text-white m-2">
				<h2 class="text-xl">Playlists</h2>
				<div class="flex flex-row overflow-x-auto overflow-y-hidden gap-2">
					{#each results.playlists as playlist}
						<Card
							id={playlist.id}
							title={playlist.name}
							subtitle={'Playlist'}
							thumbnail={playlist.thumbnail}
							type="PLAYLIST"
						/>
					{/each}
				</div>
			</div>
		</div>
	{:catch error}
		<span>{error.message}</span>
	{/await}
{/if}
