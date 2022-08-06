<script lang="ts">
	import { search } from '$lib/api/music.api';
	import Loader from '$lib/components/Loader.svelte';
	import type { MusicSearchResult } from '$lib/model/music-search-result';
	import { onMount } from 'svelte';
	import Card from '../../lib/components/Card.svelte';
	import SearchBar from '../../lib/components/SearchBar.svelte';

	let searchPromise: Promise<MusicSearchResult>;

	function submitSearch(e: any): void {
		e.preventDefault();
		searchPromise = search(e.detail);
	}

	onMount(() => {
		searchPromise = search('test');
	});
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
						background: linear-gradient(
							to top right,
							#03a9f4 35%,
							rgba(109, 5, 255, 0.6194852941) 100%
						);
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
