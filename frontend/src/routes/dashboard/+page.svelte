<script>
    import { onMount } from "svelte";
    import { goto } from "$app/navigation";
    import { Plus, Server, Play, Square, Settings } from "lucide-svelte";
    import { servers } from "$lib/stores.js";
    import CreateServerModal from "$lib/components/CreateServerModal.svelte";
    import { api } from "$lib/api";

    let isModalOpen = $state(false);

    onMount(async () => {
        try {
            // Auth check
            await api.get("/api/me");

            // Load servers
            const resServers = await api.get("/api/servers");
            servers.set(resServers.data);
        } catch (e) {
            console.error(e);
            // Redirection is handled by interceptor but we catch here to stop execution
        }
    });

    function openModal() {
        isModalOpen = true;
    }
</script>

<div class="p-8">
    <div class="mx-auto max-w-7xl">
        <!-- Header -->
        <div class="flex justify-between items-center mb-8">
            <h1 class="text-3xl font-bold text-white flex items-center gap-3">
                <Server class="text-green-500" />
                Mes Serveurs
            </h1>
            <button
                onclick={openModal}
                class="bg-green-600 hover:bg-green-500 text-white px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2 cursor-pointer"
            >
                <Plus size={20} />
                Nouveau Serveur
            </button>
        </div>

        <!-- Grid -->
        {#if $servers.length === 0}
            <div
                class="text-center py-20 bg-slate-800/50 rounded-2xl border-2 border-dashed border-slate-700"
            >
                <div
                    class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-700/50 text-slate-400 mb-4"
                >
                    <Server size={32} />
                </div>
                <h3 class="text-xl font-medium text-white mb-2">
                    Aucun serveur
                </h3>
                <p class="text-slate-400 mb-6">
                    Vous n'avez pas encore de serveur Minecraft.
                </p>
                <button
                    onclick={openModal}
                    class="mx-auto bg-slate-700 hover:bg-slate-600 text-white px-4 py-2 rounded-lg transition-colors"
                >
                    Créer mon premier serveur
                </button>
            </div>
        {:else}
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {#each $servers as server (server.id)}
                    <div
                        class="bg-slate-800 rounded-xl border border-slate-700 overflow-hidden hover:border-slate-600 transition-all hover:shadow-xl group"
                    >
                        <div class="p-6">
                            <div class="flex justify-between items-start mb-4">
                                <div>
                                    <h3
                                        class="text-lg font-bold text-white mb-1 group-hover:text-green-400 transition-colors"
                                    >
                                        {server.name}
                                    </h3>
                                    <div
                                        class="flex items-center gap-2 text-xs font-mono text-slate-400"
                                    >
                                        <span
                                            class="px-1.5 py-0.5 rounded bg-slate-700 border border-slate-600"
                                            >:{server.port}</span
                                        >
                                        <span
                                            class="px-1.5 py-0.5 rounded bg-slate-700 border border-slate-600"
                                            >{server.type}</span
                                        >
                                    </div>
                                </div>
                                <div
                                    class="px-2.5 py-1 rounded-full text-xs font-bold bg-slate-700 text-slate-400"
                                >
                                    STOPPED
                                </div>
                            </div>

                            <div
                                class="flex items-center justify-between mt-6 pt-4 border-t border-slate-700/50"
                            >
                                <div class="text-sm text-slate-500">
                                    {server.ram} MB RAM
                                </div>
                                <a
                                    href="/server/{server.id}"
                                    class="text-sm font-medium text-green-400 hover:text-green-300 flex items-center gap-1 hover:underline"
                                >
                                    <Settings size={16} />
                                    Gérer
                                </a>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>

<CreateServerModal isOpen={isModalOpen} onClose={() => (isModalOpen = false)} />
