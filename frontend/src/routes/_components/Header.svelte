<script lang="ts">
	import appIcon from '$lib/assets/icons/app_white_24dp.svg';
	import { onMount } from 'svelte';
	import { getLatency } from '../../lib/api/bot.api';
	import User from './User.svelte';

	const LATENCY_POLL_INTERVAL_MS = 5000;

	let latency: number;
	onMount(async () => {
		const res = await getLatency();
		latency = res.latency;
		setInterval(async () => {
			const res = await getLatency();
			latency = res.latency;
		}, LATENCY_POLL_INTERVAL_MS);
	});
</script>

<div class="header bg-gradient-to-br from-purple-600 to-blue-500">
	<img src={appIcon} alt="logo" class="logo" />
	<span>Discord Sound Board</span>
	<div class="flex ml-auto justify-items-end">
		<div class="latency">
			<div class="icon" />
			<div class="value">{latency} ms</div>
		</div>
		<User />
		


	</div>
</div>

<style>
	.header {
		width: 100%;
		height: 100%;
		display: flex;
		align-items: center;
		gap: 16px;
		padding: 0 16px;
		color: white;
		box-sizing: border-box;
	}

	.logo {
		height: calc(100% - 12px);
	}

	.latency {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 6px;
	}

	.icon {
		width: 8px;
		height: 8px;
		background-color: rgb(0, 199, 0);
		border-radius: 99px;
	}

	.value {
		font-family: 'Roboto-Thin', sans-serif;
		font-size: 11px;
	}
</style>
