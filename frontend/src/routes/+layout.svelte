<script lang="ts">
	import "./layout.css";
	import favicon from "$lib/assets/favicon.svg";
	import { onMount } from "svelte";
	import { goto } from "$app/navigation";
	import { page } from "$app/stores";
	import { api } from "$lib/api";

	let { children } = $props();
	let user = $state(null);
	let isMenuOpen = $state(false);

	async function checkAuth() {
		try {
			const res = await api.get("/api/me");
			user = res.data;
		} catch (e) {
			user = null;
		}
	}

	async function handleLogout() {
		try {
			await api.post("/api/logout");
			user = null;
			goto("/");
		} catch (e) {
			console.error("Logout failed", e);
		}
	}

	onMount(() => {
		checkAuth();
	});

	// Check auth on navigation
	$effect(() => {
		// Just triggering reactivity on page change if needed,
		// but explicit checks in dashboard are better.
		// We can re-check auth periodically or on focus if we wanted.
		const _ = $page.url.pathname;
		checkAuth();
	});
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
	<title>Crafteur - Gérez votre serveur Minecraft</title>
</svelte:head>

<div class="min-h-screen bg-slate-900 text-slate-100 flex flex-col font-sans">
	<!-- Navbar -->
	<nav
		class="border-b border-slate-800 bg-slate-900/50 backdrop-blur-md sticky top-0 z-50"
	>
		<div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
			<div class="flex h-16 items-center justify-between">
				<div class="flex items-center">
					<a href="/" class="flex items-center space-x-2">
						<!-- Simple Logo placeholder using text/icon -->
						<span
							class="text-2xl font-bold bg-gradient-to-r from-green-400 to-emerald-600 bg-clip-text text-transparent"
							>Crafteur</span
						>
					</a>
				</div>

				<div class="hidden md:block">
					<div class="ml-10 flex items-baseline space-x-4">
						{#if user}
							<a
								href="/dashboard"
								class="rounded-md px-3 py-2 text-sm font-medium text-slate-300 hover:bg-slate-700 hover:text-white transition-colors"
								>Dashboard</a
							>
							<button
								class="rounded-md px-3 py-2 text-sm font-medium text-red-400 hover:bg-red-500/10 hover:text-red-300 transition-colors cursor-pointer"
								>Déconnexion</button
							>
						{:else}
							<a
								href="/login"
								class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-500 transition-colors shadow-lg shadow-green-900/20"
								>Connexion</a
							>
						{/if}
					</div>
				</div>

				<!-- Mobile menu button -->
				<div class="-mr-2 flex md:hidden">
					<button
						onclick={() => (isMenuOpen = !isMenuOpen)}
						class="inline-flex items-center justify-center rounded-md bg-slate-800 p-2 text-slate-400 hover:bg-slate-700 hover:text-white focus:outline-none focus:ring-2 focus:ring-white focus:ring-offset-2 focus:ring-offset-slate-800 cursor-pointer"
					>
						<span class="sr-only">Open main menu</span>
						{#if !isMenuOpen}
							<svg
								class="block h-6 w-6"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								><path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25h16.5"
								/></svg
							>
						{:else}
							<svg
								class="block h-6 w-6"
								fill="none"
								viewBox="0 0 24 24"
								stroke-width="1.5"
								stroke="currentColor"
								><path
									stroke-linecap="round"
									stroke-linejoin="round"
									d="M6 18L18 6M6 6l12 12"
								/></svg
							>
						{/if}
					</button>
				</div>
			</div>
		</div>

		<!-- Mobile menu -->
		{#if isMenuOpen}
			<div class="md:hidden">
				<div class="space-y-1 px-2 pb-3 pt-2 sm:px-3">
					{#if user}
						<a
							href="/dashboard"
							class="block rounded-md px-3 py-2 text-base font-medium text-slate-300 hover:bg-slate-700 hover:text-white"
							>Dashboard</a
						>
						<button
							onclick={handleLogout}
							class="block w-full text-left rounded-md px-3 py-2 text-base font-medium text-red-400 hover:bg-red-500/10 hover:text-red-300 cursor-pointer"
							>Déconnexion</button
						>
					{:else}
						<a
							href="/login"
							class="block rounded-md px-3 py-2 text-base font-medium text-slate-300 hover:bg-slate-700 hover:text-white"
							>Connexion</a
						>
					{/if}
				</div>
			</div>
		{/if}
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
