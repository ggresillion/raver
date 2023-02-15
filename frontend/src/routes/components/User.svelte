<script lang="ts">
	import { onMount } from 'svelte';
	import { logout } from '$lib/api/auth.api';
	import { getUser } from '$lib/api/user.api';

	let isOpen = false;

	function toggle() {
		isOpen = !isOpen;
	}

	function signOut() {
		logout();
		window.location.replace('/login');
	}

	onMount(() => {
		window.onclick = function (event: any) {
			if (!event.target.parentElement.matches('.dropdown')) {
				isOpen = false;
			}
		};
	});
</script>

{#await getUser() then user}
	<div class="flex items-center gap-2">
		<span><span class="font-thin">Welcome, </span>{user?.username}</span>
		<div class="dropdown">
			<img
				class="rounded-full h-12 w-12 cursor-pointer"
				src={`https://cdn.discordapp.com/avatars/${user?.id}/${user?.avatar}.png`}
				alt="avatar"
				on:click={() => toggle()}
				on:keydown={() => {}}
				on:keyup={() => {}}
				on:keypress={() => {}}
			/>
			<div
				class="group z-10 w-44 bg-white rounded divide-y divide-gray-100 shadow dark:bg-gray-700 block absolute right-0 top-16"
				class:invisible={!isOpen}
			>
				<ul class="py-1 text-sm text-gray-700 dark:text-gray-200" aria-labelledby="dropdownDefault">
					<li>
						<a
							href={'#'}
							class="block py-2 px-4 hover:bg-gray-100 dark:hover:bg-gray-600 dark:hover:text-white"
							on:click={() => signOut()}>Sign out</a
						>
					</li>
				</ul>
			</div>
		</div>
	</div>
{/await}
