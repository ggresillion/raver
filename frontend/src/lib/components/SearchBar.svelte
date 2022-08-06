<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	let input: string;

	const dispatch = createEventDispatcher();

	function submit(e: SubmitEvent): void {
		e.preventDefault();
		dispatch('submit', input);
	}
</script>

<form on:submit={submit}>
	<div class="relative w-96">
		<div class="flex absolute inset-y-0 left-0 items-center pl-3 pointer-events-none">
			<svg
				aria-hidden="true"
				class="w-5 h-5 text-gray-500 dark:text-gray-400"
				fill="none"
				stroke="currentColor"
				viewBox="0 0 24 24"
				xmlns="http://www.w3.org/2000/svg"
				><path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
				/></svg
			>
		</div>
		<input
			type="search"
			class="block p-4 pl-10 w-full text-sm text-gray-900 bg-[#414141] rounded-lg border border-gray-600 focus:ring-blue-500 focus:border-blue-500"
			placeholder="Search Songs, Artists..."
			required
			bind:value={input}
		/>
		<button
			type="submit"
			class="text-white absolute right-2.5 bottom-2.5 bg-gradient-to-br from-purple-600 to-blue-500 hover:bg-gradient-to-bl focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2"
			>Search</button
		>
	</div>
</form>
