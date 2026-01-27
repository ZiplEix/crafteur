<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import {
        Search,
        Download,
        Check,
        ExternalLink,
        Loader2,
        Info,
    } from "lucide-svelte";
    import type {
        ModrinthProject,
        ModrinthSearchResponse,
    } from "$lib/types/modrinth";

    export let serverId: string;

    let query = "";
    let loading = false;
    let hits: ModrinthProject[] = [];
    let debounceTimer: any;
    let installing: Record<string, boolean> = {};
    let installed: Record<string, boolean> = {}; // Track session-installed items

    async function search() {
        loading = true;
        try {
            const res = await api.get(`/api/modrinth/search`, {
                params: {
                    q: query,
                    serverId: serverId,
                    limit: 20,
                },
            });
            const data: ModrinthSearchResponse = res.data;
            hits = data.hits;
        } catch (e) {
            console.error("Search failed", e);
        } finally {
            loading = false;
        }
    }

    function handleInput() {
        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            search();
        }, 500);
    }

    async function install(project: ModrinthProject) {
        if (installing[project.project_id]) return;

        installing[project.project_id] = true;
        try {
            await api.post(`/api/modrinth/install`, {
                serverId: serverId,
                projectId: project.project_id,
            });
            installed[project.project_id] = true;
            // Maybe show a toast via a global store or event?
            // For now relies on button state change
        } catch (e: any) {
            console.error("Install failed", e);
            alert(
                "Erreur d'installation : " +
                    (e.response?.data?.error || e.message),
            );
        } finally {
            installing[project.project_id] = false;
        }
    }

    function formatDownloads(count: number): string {
        if (count >= 1000000) {
            return (count / 1000000).toFixed(1) + "M";
        }
        if (count >= 1000) {
            return (count / 1000).toFixed(1) + "k";
        }
        return count.toString();
    }

    onMount(() => {
        search(); // Initial load
    });
</script>

<div class="space-y-6">
    <!-- Search Bar -->
    <div class="relative">
        <div
            class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-gray-400"
        >
            <Search size={20} />
        </div>
        <input
            type="text"
            bind:value={query}
            on:input={handleInput}
            placeholder="Rechercher un mod (ex: JEI, Sodium, Create...)"
            class="block w-full pl-10 pr-3 py-3 border border-gray-700 rounded-lg leading-5 bg-gray-900/50 text-gray-300 placeholder-gray-500 focus:outline-none focus:bg-gray-900 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm transition-colors"
        />
    </div>

    <!-- Info Info -->
    <div
        class="flex items-center gap-2 text-xs text-blue-400 bg-blue-500/10 p-3 rounded-lg border border-blue-500/20"
    >
        <Info size={16} />
        <span
            >Filtre automatique appliqué : Compatibilité vérifiée avec la
            version et le type de votre serveur.</span
        >
    </div>

    <!-- Grid -->
    {#if loading && hits.length === 0}
        <div class="flex justify-center p-12">
            <Loader2 size={32} class="animate-spin text-blue-500" />
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each hits as project}
                <div
                    class="bg-gray-800/40 border border-gray-700/50 rounded-xl p-4 hover:bg-gray-800/60 transition-all flex flex-col gap-3 group"
                >
                    <div class="flex gap-4">
                        {#if project.icon_url}
                            <img
                                src={project.icon_url}
                                alt={project.title}
                                class="w-12 h-12 rounded-lg bg-gray-900 object-cover"
                            />
                        {:else}
                            <div
                                class="w-12 h-12 rounded-lg bg-gray-700 flex items-center justify-center text-gray-500"
                            >
                                <Search size={20} />
                            </div>
                        {/if}
                        <div class="flex-1 min-w-0">
                            <h3
                                class="text-sm font-bold text-gray-200 truncate pr-2"
                            >
                                {project.title}
                            </h3>
                            <div
                                class="flex items-center gap-2 text-xs text-gray-500"
                            >
                                <span class="truncate max-w-[80px]"
                                    >{project.author}</span
                                >
                                <span class="text-gray-600">•</span>
                                <span class="flex items-center gap-1">
                                    <Download size={10} />
                                    {formatDownloads(project.downloads)}
                                </span>
                            </div>
                            <div class="flex gap-1 mt-1">
                                {#if project.client_side === "required" || project.client_side === "optional"}
                                    <span
                                        class="text-[10px] bg-emerald-500/10 text-emerald-400 px-1.5 py-0.5 rounded border border-emerald-500/20"
                                        >Client</span
                                    >
                                {/if}
                                {#if project.server_side === "required" || project.server_side === "optional"}
                                    <span
                                        class="text-[10px] bg-indigo-500/10 text-indigo-400 px-1.5 py-0.5 rounded border border-indigo-500/20"
                                        >Server</span
                                    >
                                {/if}
                            </div>
                        </div>
                    </div>

                    <p class="text-xs text-gray-400 line-clamp-2 h-8">
                        {project.description}
                    </p>

                    <div class="mt-auto pt-2 flex items-center justify-between">
                        <a
                            href={`https://modrinth.com/project/${project.slug}`}
                            target="_blank"
                            rel="noopener noreferrer"
                            class="text-xs text-gray-500 hover:text-gray-300 flex items-center gap-1 transition-colors"
                        >
                            <ExternalLink size={12} />
                            Modrinth
                        </a>

                        {#if installed[project.project_id]}
                            <button
                                disabled
                                class="flex items-center gap-2 px-3 py-1.5 rounded-md text-xs font-medium bg-gray-700 text-gray-400 cursor-not-allowed"
                            >
                                <Check size={14} />
                                Installé
                            </button>
                        {:else}
                            <button
                                on:click={() => install(project)}
                                disabled={installing[project.project_id]}
                                class="flex items-center gap-2 px-3 py-1.5 rounded-md text-xs font-medium bg-green-600 hover:bg-green-500 text-white transition-colors disabled:opacity-50 disabled:cursor-wait shadow-lg shadow-green-900/10"
                            >
                                {#if installing[project.project_id]}
                                    <Loader2 size={14} class="animate-spin" />
                                    Installation...
                                {:else}
                                    <Download size={14} />
                                    Installer
                                {/if}
                            </button>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>
