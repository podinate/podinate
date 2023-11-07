<!-- YOU CAN DELETE EVERYTHING IN THIS PAGE -->

<div class="container h-full mx-auto flex justify-center items-center">
	<div class="space-y-10 text-center flex flex-col items-center">
		<h2 class="h2">Log In</h2>
		<!-- Animated Logo -->
		<figure>
			<!-- <section class="img-bg">
			Heyyyyy
			</section> -->
			<div class="card p-4 space-y-2 grid-cols-4">
				Log in with...<br />
				<button type="button" class="btn btn-lg variant-filled flex-1 w-64">Github</button><br />
				<button type="button" class="btn btn-lg variant-filled flex-1 w-64">Gitlab</button><br />
				Or...<br />
				<button type="button" class="btn btn-lg variant-filled flex-1 w-64">Create Account</button>
			</div>
		</figure>
		<!-- / -->
		<div class="flex justify-center space-x-2">
			
		</div>
		<div class="space-y-2">
			<p>Try editing these following:</p>
			<p><code class="code">/src/routes/+layout.svelte</code></p>
			<p><code class="code">/src/routes/+page.svelte</code></p>
		</div>
	</div>
</div>

<style lang="postcss">
	figure {
		@apply flex relative flex-col;
	}
	figure svg,
	.img-bg {
		@apply w-64 h-64 md:w-80 md:h-80;
	}
	.img-bg {
		@apply absolute z-[-1] rounded-full blur-[50px] transition-all;
		animation: pulse 5s cubic-bezier(0, 0, 0, 0.5) infinite,
			glow 5s linear infinite;
	}
	@keyframes glow {
		0% {
			@apply bg-primary-400/50;
		}
		33% {
			@apply bg-secondary-400/50;
		}
		66% {
			@apply bg-tertiary-400/50;
		}
		100% {
			@apply bg-primary-400/50;
		}
	}
	@keyframes pulse {
		50% {
			transform: scale(1.5);
		}
	}
</style>

<script>
	
	import { onMount } from 'svelte';
	import {Account} from '../../lib/api_client';
	onMount(() => {
		var benis = 'benis';
		console.log('Hello World!');
	});
</script>