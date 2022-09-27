<script lang="ts">
	import { join } from '$lib/api/bot.api.js';
	import {
		getPlayerState,
		subscribeToPlayerState,
		subscribeToProgress
	} from '$lib/api/music.api';
	import { pause, play, skip } from '$lib/api/music.api.js';
	import Loader from '$lib/components/Loader.svelte';
	import { PlayerStatus } from '$lib/model/player-status';
	import { selectedGuildId } from '$lib/stores/guild.store';
	import { playerState } from '$lib/stores/player-state.store';
	import Soundwave from './Soundwave.svelte';

	let progress = 0;

	function formatMillisToMinutesAndSeconds(millis: number): string {
		const time = {
			minutes: Math.floor(millis / 60000),
			seconds: ((millis % 60000) / 1000).toFixed(0)
		};
		if (!millis) return '-:-';
		const seconds = (parseInt(time.seconds) < 10 ? '0' : '') + time.seconds;
		return time.minutes + ':' + seconds;
	}

	function progressToPercent(actual: number, total: number): number {
		if (!actual || !total || isNaN(actual) || isNaN(total)) return 0;
		return Math.round((actual / total) * 100);
	}

	playerState.subscribe((p) => {
		if (!p || !p.progress) return;
		if (!p.playlist || p.playlist.length < 1) return;
		progress = progressToPercent(p.progress, p.playlist[0].duration);
	});

	selectedGuildId.subscribe(async (val) => {
		if (!val || !$selectedGuildId) return;
		playerState.set(await getPlayerState($selectedGuildId));
		await subscribeToPlayerState(val);
		await subscribeToProgress(val);
	});
</script>

<div class="player">
	{#if !$selectedGuildId}
		<div />
	{:else if !$playerState}
		<Loader />
	{:else}
		<div class="track">
			{#if $playerState.playlist.length > 0}
				<img class="track-thumbnail" src={$playerState.playlist[0].thumbnail} alt="thumbnail" />
				<div class="track-info">
					<span class="track-name">{$playerState.playlist[0].title}</span>
					<span class="track-artist"
						>{$playerState.playlist[0].artists.map((a) => a.name).join(', ')}</span
					>
				</div>
				<Soundwave play={$playerState.status === PlayerStatus.PLAYING} />
			{/if}
		</div>

		<div class="controls">
			<div>
				<button type="button" style="visibility: hidden" />
				{#if $playerState.status === PlayerStatus.IDLE || $playerState.status === PlayerStatus.PAUSED}
					<button
						type="button"
						class="play"
						on:click={() => !$selectedGuildId || play($selectedGuildId)}
					/>
				{:else if $playerState.status === PlayerStatus.PLAYING}
					<button
						type="button"
						class="pause"
						on:click={() => !$selectedGuildId || pause($selectedGuildId)}
					/>
				{:else if $playerState.status === PlayerStatus.BUFFERING}
					<button type="button" class="buffering" disabled />
				{:else}
					<button type="button" class="play" disabled />
				{/if}
				<button
					type="button"
					class="skip-next"
					on:click={() => !$selectedGuildId || skip($selectedGuildId)}
					disabled={!$playerState || $playerState.playlist.length <= 0}
				/>
			</div>

			<div>
				{#if $playerState.playlist.length > 0 && $playerState.progress}
					<div>{formatMillisToMinutesAndSeconds($playerState.progress)}</div>
					<div class="progress-bar">
						<div class="progress-bar">
							<input
								id="disabled-range"
								type="range"
								value={progress}
								class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
								disabled
							/>
						</div>
					</div>
				{:else}
					<div>-:-</div>
					<div class="progress-bar">
						<input
							id="disabled-range"
							type="range"
							value="0"
							class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer dark:bg-gray-700"
							disabled
						/>
					</div>
				{/if}
				{#if $playerState.playlist.length > 0}
					<span>{formatMillisToMinutesAndSeconds($playerState.playlist[0].duration)}</span>
				{:else}
					<span>-:-</span>
				{/if}
			</div>
		</div>

		<div class="extra">
			<button
				type="button"
				class="text-white bg-gradient-to-br from-purple-600 to-blue-500 hover:bg-gradient-to-bl focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-5 py-2.5 text-center mr-2 mb-2"
				title="Join voice channel"
				disabled={!$selectedGuildId}
				on:click={() => !$selectedGuildId || join($selectedGuildId)}>Join my channel</button
			>
		</div>
	{/if}
</div>

<style lang="scss">
	.player {
		position: fixed;
		bottom: 0;
		left: 0;
		height: 100px;
		width: 100vw;
		padding: 6px;
		background-color: rgb(53, 53, 53);
		filter: drop-shadow(3px 5px 2px rgb(0 0 0 / 0.4));
		color: white;
		box-sizing: border-box;

		display: grid;
		grid-template-columns: repeat(3, 1fr);
		grid-gap: 10px;

		.track {
			display: flex;
			align-items: center;
			gap: 26px;

			.track-thumbnail {
				height: 80px;
				filter: none;
			}

			.track-info {
				display: flex;
				flex-direction: column;
				justify-content: space-between;
				padding: 4px 0;

				.track-name {
					font-family: 'Roboto-Thin', serif;
				}

				.track-artist {
					font-family: 'Roboto', serif;
				}
			}
		}

		.controls {
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
			gap: 6px;

			& > div {
				display: flex;
				align-items: center;
				justify-content: center;
				gap: 6px;
				width: 100%;
			}

			.progress-bar {
				flex: 1;
			}

			button {
				cursor: pointer;
				user-select: none;
				background-color: white;
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

				&.pause {
					background: url('$lib/assets/icons/pause.svg');
				}

				&.play {
					background: url('$lib/assets/icons/play.svg');
				}

				&.buffering {
					background: url('$lib/assets/icons/buffering.svg');
					animation: rotation 2s infinite linear;

					@keyframes rotation {
						from {
							transform: rotate(0deg);
						}
						to {
							transform: rotate(359deg);
						}
					}
				}

				&.skip-next {
					background: url('$lib/assets/icons/skip_next.svg');
					width: 32px;
					height: 32px;
				}
			}
		}

		.extra {
			display: flex;
			justify-content: flex-end;
			align-items: center;
		}

		button {
			&:disabled {
				opacity: 0.2;
				cursor: default;

				&:hover {
					transform: none;
				}
			}

			&:hover {
				transform: scale(1.1);
			}
		}
	}

	//.slider {
	//  width: 100%;
	//  height: 10px;
	//  border: 1px rgb(27, 27, 27) solid;
	//  box-shadow: 0px 0px 25px black;
	//  border-radius: 99px;
	//  user-select: none;
	//  background: linear-gradient(90deg, #03a9f4 35%, rgba(109, 5, 255, 0.6194852941) 100%);
	//
	//  &:hover {
	//    .thumb {
	//      transform: scale(1.1);
	//    }
	//  }
	//
	//  .thumb {
	//    height: 12px;
	//    width: 12px;
	//    margin-top: -1px;
	//    background-color: white;
	//    border-radius: 99px;
	//    cursor: pointer;
	//  }
	//}
</style>
