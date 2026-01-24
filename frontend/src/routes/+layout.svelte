<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { Hammer, Github, LogIn } from "lucide-svelte";

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Crafteur - Gérez votre serveur Minecraft</title>
</svelte:head>

<div class="min-h-screen bg-slate-900 text-slate-100 flex flex-col font-sans">
	<!-- Navbar -->
	<nav
		class="sticky top-0 z-50 border-b border-slate-800 bg-slate-900/80 backdrop-blur-md"
	>
		<div
			class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8"
		>
			<!-- Logo -->
			<a
				href="/"
				class="flex items-center gap-2 transition hover:opacity-80"
			>
				<Hammer class="h-6 w-6 text-green-500" />
				<span class="text-xl font-bold tracking-tight text-white"
					>Crafteur</span
				>
			</a>

			<!-- Right Actions -->
			<a
				href="/login"
				class="group inline-flex items-center gap-2 rounded-md bg-slate-800 px-4 py-2 text-sm font-medium text-slate-200 transition hover:bg-slate-700 hover:text-white"
			>
				<LogIn class="h-4 w-4 transition group-hover:translate-x-0.5" />
				<span>Connexion</span>
			</a>
		</div>
	</nav>

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
