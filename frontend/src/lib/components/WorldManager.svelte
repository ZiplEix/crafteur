<script lang="ts">
    import { onMount } from "svelte";
    import { Globe, Trash, Check, Plus, AlertTriangle } from "lucide-svelte";
    import { api } from "$lib/api";

    export let serverId: string;

    interface WorldEntry {
        name: string;
        is_active: boolean;
        size: number;
    }

    let worlds: WorldEntry[] = [];
    let loading = true;
    let error: string | null = null;
    let showCreateModal = false;
    let newWorldName = "";
    let creating = false;

    // Load worlds
    async function fetchWorlds() {
        loading = true;
        try {
            const res = await api.get(`/api/servers/${serverId}/worlds`);
            worlds = res.data;
            error = null;
        } catch (e: any) {
            console.error("Failed to fetch worlds", e);
            error = e.message;
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        fetchWorlds();
    });

    // Create World
    async function createWorld() {
        if (!newWorldName.trim()) return;

        // Simple client request validation
        if (!/^[a-zA-Z0-9_-]+$/.test(newWorldName)) {
            alert(
                "Nom invalide (alphanurérique, tirets et underscores uniquement)",
            );
            return;
        }

        creating = true;
        try {
            await api.post(`/api/servers/${serverId}/worlds`, {
                name: newWorldName,
            });
            await fetchWorlds();
            showCreateModal = false;
            newWorldName = "";
        } catch (e: any) {
            alert(
                "Erreur lors de la création: " +
                    (e.response?.data?.error || e.message),
            );
        } finally {
            creating = false;
        }
    }

    // Activate World
    async function activateWorld(name: string) {
        if (
            !confirm(
                `Voulez-vous vraiment activer le monde "${name}" ?\nLe serveur devra être redémarré.`,
            )
        ) {
            return;
        }

        try {
            await api.post(`/api/servers/${serverId}/worlds/${name}/activate`);
            await fetchWorlds(); // Refresh to see update (active flag)
            alert(
                "Monde modifié ! Redémarrez le serveur pour appliquer les changements.",
            );
        } catch (e: any) {
            alert(
                "Erreur lors de l'activation: " +
                    (e.response?.data?.error || e.message),
            );
        }
    }

    // Delete World
    async function deleteWorld(name: string) {
        if (
            !confirm(
                `Supprimer dÃ©finitivement le monde "${name}" ?\nCette action est irréversible.`,
            )
        ) {
            return;
        }

        try {
            await api.delete(`/api/servers/${serverId}/worlds/${name}`);
            await fetchWorlds();
        } catch (e: any) {
            alert(
                "Erreur lors de la suppression: " +
                    (e.response?.data?.error || e.message),
            );
        }
    }

    function formatBytes(bytes: number): string {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }
</script>

<div class="space-y-6">
    <!-- Header -->
    <div
        class="flex justify-between items-center bg-slate-800/50 p-4 rounded-xl border border-gray-700"
    >
        <div>
            <h2 class="text-xl font-bold text-white flex items-center gap-2">
                <Globe class="text-blue-400" />
                Gestion des Mondes
            </h2>
            <p class="text-gray-400 text-sm mt-1">
                Gérez vos mondes Minecraft. Le monde actif est celui chargé au
                démarrage.
            </p>
        </div>
        <button
            on:click={() => (showCreateModal = true)}
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium flex items-center gap-2 transition-colors"
        >
            <Plus size={18} />
            Nouveau Monde
        </button>
    </div>

    <!-- Error/Loading -->
    {#if loading}
        <div class="text-center py-10 text-gray-400">
            Chargement des mondes...
        </div>
    {:else if error}
        <div
            class="bg-red-500/10 border border-red-500/20 text-red-400 p-4 rounded-lg flex items-center gap-2"
        >
            <AlertTriangle size={20} />
            {error}
        </div>
    {:else if worlds.length === 0}
        <div class="text-center py-10 text-gray-500">
            Aucun monde trouvé. Créez-en un nouveau !
        </div>
    {:else}
        <!-- Grid -->
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {#each worlds as world}
                <div
                    class="relative group bg-slate-800 rounded-xl border transition-all duration-200 p-5 flex flex-col gap-4
                    {world.is_active
                        ? 'border-green-500/50 shadow-[0_0_15px_rgba(34,197,94,0.1)]'
                        : 'border-gray-700 hover:border-gray-500 hover:bg-slate-750'}"
                >
                    <!-- Header Card -->
                    <div class="flex justify-between items-start">
                        <div class="flex items-center gap-3">
                            <div
                                class="p-2 rounded-lg {world.is_active
                                    ? 'bg-green-500/20 text-green-400'
                                    : 'bg-gray-700/50 text-gray-400'}"
                            >
                                <Globe size={24} />
                            </div>
                            <div>
                                <h3
                                    class="font-semibold text-white text-lg leading-tight"
                                >
                                    {world.name}
                                </h3>
                                <div
                                    class="text-xs text-gray-500 font-mono mt-1"
                                >
                                    {formatBytes(world.size)}
                                </div>
                            </div>
                        </div>
                        {#if world.is_active}
                            <span
                                class="px-2 py-0.5 rounded text-xs font-bold bg-green-500 text-black uppercase tracking-wide"
                            >
                                Actif
                            </span>
                        {/if}
                    </div>

                    <!-- Actions -->
                    <div
                        class="pt-2 mt-auto border-t border-gray-700/50 flex gap-2 justify-end"
                    >
                        {#if !world.is_active}
                            <button
                                on:click={() => activateWorld(world.name)}
                                class="flex-1 bg-slate-700 hover:bg-blue-600 text-white py-2 px-3 rounded-lg text-sm font-medium transition-colors flex items-center justify-center gap-2"
                                title="Activer ce monde"
                            >
                                <Check size={16} />
                                Activer
                            </button>
                            <button
                                on:click={() => deleteWorld(world.name)}
                                class="bg-red-500/10 hover:bg-red-600 text-red-400 hover:text-white p-2 rounded-lg transition-colors shrink-0"
                                title="Supprimer définitivement"
                            >
                                <Trash size={18} />
                            </button>
                        {:else}
                            <div
                                class="text-sm text-green-400/80 italic w-full text-center py-1.5 cursor-default"
                            >
                                Monde actuellement chargé
                            </div>
                        {/if}
                    </div>
                </div>
            {/each}
        </div>
    {/if}
</div>

<!-- Modal Creation -->
{#if showCreateModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center px-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-gray-700 rounded-xl p-6 w-full max-w-md shadow-2xl transform transition-all"
        >
            <h3 class="text-xl font-bold text-white mb-4">
                Créer un nouveau monde
            </h3>

            <div class="space-y-4">
                <div>
                    <label
                        for="worldName"
                        class="block text-sm font-medium text-gray-400 mb-1"
                        >Nom du monde</label
                    >
                    <input
                        id="worldName"
                        type="text"
                        bind:value={newWorldName}
                        placeholder="Ex: survival_2024"
                        class="w-full bg-slate-950 border border-gray-700 rounded-lg p-3 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    />
                    <p class="text-xs text-gray-500 mt-1">
                        Caractères autorisés : lettres, chiffres, tirets,
                        underscores.
                    </p>
                </div>

                <div class="flex gap-3 justify-end pt-2">
                    <button
                        on:click={() => {
                            showCreateModal = false;
                            newWorldName = "";
                        }}
                        class="px-4 py-2 rounded-lg text-gray-300 hover:text-white hover:bg-slate-800 transition-colors"
                        disabled={creating}
                    >
                        Annuler
                    </button>
                    <button
                        on:click={createWorld}
                        disabled={!newWorldName || creating}
                        class="bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed text-white px-4 py-2 rounded-lg font-medium transition-colors flex items-center gap-2"
                    >
                        {#if creating}
                            <div
                                class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"
                            ></div>
                        {/if}
                        Créer
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}
