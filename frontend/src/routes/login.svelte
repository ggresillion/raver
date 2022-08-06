<script lang="ts">
	import Particles from 'svelte-particles';
	import { config } from '../lib/api/config.api';
	import { loadFull } from 'tsparticles';

	const particlesConfig: any = {
		background: {
			color: {
				value: '#232741'
			},
			position: '50% 50%',
			repeat: 'no-repeat',
			size: '20%'
		},
		fullScreen: {
			zIndex: 1
		},
		interactivity: {
			events: {
				onClick: {
					enable: true,
					mode: 'repulse'
				},
				onHover: {
					enable: true,
					mode: 'bubble'
				}
			},
			modes: {
				bubble: {
					distance: 250,
					duration: 2,
					opacity: 0,
					size: 0,
					divs: {
						distance: 200,
						duration: 0.4,
						mix: false,
						selectors: []
					}
				},
				grab: {
					distance: 400
				},
				repulse: {
					distance: 400,
					divs: {
						distance: 200,
						duration: 0.4,
						factor: 100,
						speed: 1,
						maxSpeed: 50,
						easing: 'ease-out-quad',
						selectors: []
					}
				}
			}
		},
		particles: {
			color: {
				value: '#ffffff'
			},
			links: {
				color: {
					value: '#ffffff'
				},
				distance: 150,
				opacity: 0.4
			},
			move: {
				attract: {
					rotate: {
						x: 600,
						y: 600
					}
				},
				enable: true,
				outModes: {
					bottom: 'out',
					left: 'out',
					right: 'out',
					top: 'out'
				},
				random: true,
				speed: 1
			},
			number: {
				density: {
					enable: true
				},
				value: 160
			},
			opacity: {
				random: {
					enable: true
				},
				value: {
					min: 0,
					max: 1
				},
				animation: {
					enable: true,
					speed: 1,
					minimumValue: 0
				}
			},
			size: {
				random: {
					enable: true
				},
				value: {
					min: 1,
					max: 3
				},
				animation: {
					speed: 4,
					minimumValue: 0.3
				}
			}
		}
	};

	async function particlesInit(engine: any) {
		await loadFull(engine);
	}

	function login() {
		window.location.href = `${config.apiUrl}/auth/login`;
	}
</script>

<Particles id="tsparticles" options={particlesConfig} {particlesInit} />

<div class="fixed z-10 flex items-center justify-center w-screen h-screen">
	<div
		class="p-6 max-w-sm bg-white rounded-lg border border-gray-200 shadow-md dark:bg-gray-800 dark:border-gray-700"
	>
		<a href="#">
			<h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">Login</h5>
		</a>
		<p class="mb-3 font-normal text-gray-700 dark:text-gray-400">
			We need you to login with Discord
		</p>
		<button
			type="submit"
			class="text-white right-2.5 bottom-2.5 bg-gradient-to-br from-purple-600 to-blue-500 hover:bg-gradient-to-bl focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2"
			on:click={login}>Login</button
		>
	</div>
</div>
