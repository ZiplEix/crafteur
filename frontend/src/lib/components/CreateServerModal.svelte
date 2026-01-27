<script lang="ts">
    import { X, Upload } from "lucide-svelte";
    import { servers } from "$lib/stores";
    import { versions, loadVersions } from "$lib/stores/meta";
    import { api } from "$lib/api";

    interface Props {
        isOpen: boolean;
        onClose: () => void;
        onServerCreated?: () => void;
    }

    let { isOpen, onClose, onServerCreated = () => {} }: Props = $props();

    $effect(() => {
        if (isOpen) {
            loadVersions();
        }
    });

    let name = $state("");
    let port = $state(25565);
    let ram = $state(2048);
    let version = $state("1.21.4");
    let type = $state("vanilla");
    let importFiles: FileList | null = $state(null);
    let loading = $state(false);

    let error: string | null = $state(null);

    async function handleSubmit(e: Event) {
        e.preventDefault();
        loading = true;
        error = null;

        try {
            const formData = new FormData();
            formData.append("name", name);
            formData.append("port", port.toString());
            formData.append("ram", ram.toString());
            formData.append("type", type);
            formData.append("version", version);

            if (importFiles && importFiles.length > 0) {
                formData.append("import", importFiles[0]);
            }

            // We use axios (via api wrapper) which might try to set Content-Type if simply passed an object.
            // But usually passing FormData lets the browser set Content-Type with boundary.
            // checking api wrapper implementation is not possible right now but standard axios handles FormData correctly.
            const res = await api.post("/api/servers", formData);

            const newServer = res.data;

            // Update store
            servers.update((s) => [...s, newServer]);

            // Reset form
            name = "";
            port = 25565;
            ram = 2048;
            importFiles = null;

            onServerCreated();
            onClose();
        } catch (err: any) {
            if (err.response && err.response.data && err.response.data.error) {
                error = err.response.data.error;
            } else if (err instanceof Error) {
                error = err.message;
            } else {
                error = "Une erreur inconnue est survenue";
            }
        } finally {
            loading = false;
        }
    }
</script>

{#if isOpen}
    <div
        class="fixed inset-0 z-[100] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="w-full max-w-md bg-slate-800 rounded-xl border border-slate-700 shadow-2xl overflow-hidden"
        >
            <!-- Header -->
            <div
                class="px-6 py-4 border-b border-slate-700 flex justify-between items-center bg-slate-900/50"
            >
                <h2 class="text-xl font-bold text-white">Nouveau Serveur</h2>
                <button
                    onclick={onClose}
                    class="text-slate-400 hover:text-white transition-colors cursor-pointer"
                >
                    <X size={20} />
                </button>
            </div>

            <!-- Body -->
            <form onsubmit={handleSubmit} class="p-6 space-y-4">
                {#if error}
                    <div
                        class="bg-red-500/10 border border-red-500/20 text-red-400 px-4 py-3 rounded-lg text-sm"
                    >
                        {error}
                    </div>
                {/if}

                <div class="space-y-2">
                    <label
                        for="name"
                        class="block text-sm font-medium text-slate-300"
                        >Nom du serveur</label
                    >
                    <input
                        type="text"
                        id="name"
                        bind:value={name}
                        required
                        class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent text-white placeholder-slate-500 outline-none transition-all"
                        placeholder="Mon Serveur Survie"
                    />
                </div>

                <div class="space-y-2">
                    <label
                        for="type"
                        class="block text-sm font-medium text-slate-300"
                        >Type de serveur</label
                    >
                    <div class="space-y-2">
                        <select
                            id="type"
                            bind:value={type}
                            class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent text-white outline-none transition-all"
                        >
                            <option value="vanilla">Vanilla</option>
                            <option value="fabric">Fabric</option>
                            <option value="paper">Paper (Plugins)</option>
                        </select>
                        {#if type === "fabric"}
                            <p
                                class="text-xs text-blue-400 bg-blue-500/10 px-3 py-2 rounded-lg border border-blue-500/20"
                            >
                                L'installation Fabric inclura automatiquement le
                                Loader le plus récent.
                            </p>
                        {:else if type === "paper"}
                            <p
                                class="text-xs text-amber-400 bg-amber-500/10 px-3 py-2 rounded-lg border border-amber-500/20"
                            >
                                Recommandé pour les serveurs publics. Supporte
                                les plugins Spigot/Bukkit. Ne supporte pas les
                                mods Fabric/Forge.
                            </p>
                        {/if}
                    </div>
                </div>

                <div class="space-y-2">
                    <label
                        for="version"
                        class="block text-sm font-medium text-slate-300"
                        >Version Minecraft</label
                    >
                    <select
                        id="version"
                        bind:value={version}
                        class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent text-white outline-none transition-all"
                    >
                        {#each $versions as v}
                            <option value={v.id}
                                >{v.id}
                                {v.type !== "release"
                                    ? `(${v.type})`
                                    : ""}</option
                            >
                        {/each}
                    </select>
                </div>

                <div class="space-y-2">
                    <label
                        for="import"
                        class="block text-sm font-medium text-slate-300"
                        >Importer des données (Optionnel)</label
                    >
                    <div class="relative">
                        <input
                            type="file"
                            id="import"
                            accept=".zip"
                            bind:files={importFiles}
                            class="hidden"
                        />
                        <label
                            for="import"
                            class="flex items-center justify-center w-full px-3 py-2 bg-slate-900 border border-dashed border-slate-600 rounded-lg cursor-pointer hover:border-green-500 hover:text-green-500 text-slate-400 transition-all gap-2"
                        >
                            <Upload size={16} />
                            <span class="truncate">
                                {importFiles && importFiles.length > 0
                                    ? importFiles[0].name
                                    : "Sélectionner un fichier .zip"}
                            </span>
                        </label>
                    </div>
                </div>

                <div class="grid grid-cols-2 gap-4">
                    <div class="space-y-2">
                        <label
                            for="port"
                            class="block text-sm font-medium text-slate-300"
                            >Port</label
                        >
                        <input
                            type="number"
                            id="port"
                            bind:value={port}
                            class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent text-white outline-none transition-all"
                        />
                    </div>
                    <div class="space-y-2">
                        <label
                            for="ram"
                            class="block text-sm font-medium text-slate-300"
                            >RAM (MB)</label
                        >
                        <input
                            type="number"
                            id="ram"
                            bind:value={ram}
                            step="512"
                            class="w-full px-3 py-2 bg-slate-900 border border-slate-700 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent text-white outline-none transition-all"
                        />
                    </div>
                </div>

                <div class="pt-2">
                    <button
                        type="submit"
                        disabled={loading}
                        class="w-full bg-green-600 hover:bg-green-500 text-white font-semibold py-2.5 rounded-lg transition-all disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                        {#if loading}
                            <div
                                class="w-5 h-5 border-2 border-white/30 border-t-white rounded-full animate-spin"
                            ></div>
                            {importFiles && importFiles.length > 0
                                ? "Création et Import..."
                                : "Création..."}
                        {:else}
                            Créer le serveur
                        {/if}
                    </button>
                </div>
            </form>
        </div>
    </div>
{/if}
