<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Crafteur - Gérez votre serveur Minecraft</title>
</svelte:head>

<div class="min-h-screen bg-slate-900 text-slate-100 flex flex-col font-sans">
	<!-- Main Content -->
	<main class="flex-1">
		{@render children()}
	</main>

	<!-- Footer -->
	<footer class="border-t border-slate-800 bg-slate-900 py-8">
		<div class="mx-auto max-w-7xl px-4 text-center sm:px-6 lg:px-8">
			<p class="text-sm text-slate-500">
				&copy; {new Date().getFullYear()} Crafteur. Powered by
				<span class="text-slate-400">Go</span>
				& <span class="text-slate-400">SvelteKit</span>.
			</p>
		</div>
	</footer>
</div>
