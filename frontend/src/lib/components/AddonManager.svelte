<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import {
        Package,
        Zap,
        Box,
        Upload,
        Trash2,
        File as FileIcon,
        Loader2,
    } from "lucide-svelte";
    import type { AddonFile, AddonType } from "$lib/types/addons";

    export let serverId: string;

    let activeTab: AddonType = "mods";
    let files: AddonFile[] = [];
    let loading = false;
    let uploading = false;
    let error: string | null = null;
    let fileInput: HTMLInputElement;

    const tabs: { id: AddonType; label: string; icon: any; accept: string }[] =
        [
            { id: "mods", label: "Mods", icon: Package, accept: ".jar" },
            { id: "plugins", label: "Plugins", icon: Zap, accept: ".jar" },
            { id: "datapacks", label: "Datapacks", icon: Box, accept: ".zip" },
        ];

    async function loadAddons() {
        loading = true;
        error = null;
        try {
            const res = await api.get(
                `/api/servers/${serverId}/addons/${activeTab}`,
            );
            files = res.data;
        } catch (e: any) {
            console.error("Failed to load addons", e);
            error = "Impossible de charger les extensions.";
        } finally {
            loading = false;
        }
    }

    async function handleUpload(event: Event) {
        const target = event.target as HTMLInputElement;
        if (!target.files || target.files.length === 0) return;

        uploading = true;
        const formFiles = target.files;

        const formData = new FormData();
        // Append all files with the key 'files'
        for (let i = 0; i < formFiles.length; i++) {
            formData.append("files", formFiles[i]);
        }

        try {
            await api.post(
                `/api/servers/${serverId}/addons/${activeTab}`,
                formData,
                {
                    headers: {
                        "Content-Type": "multipart/form-data",
                    },
                },
            );

            alert(`${formFiles.length} fichiers ajoutés avec succès !`);
            await loadAddons();
        } catch (e: any) {
            console.error("Upload failed", e);
            alert(
                "Erreur lors de l'upload: " +
                    (e.response?.data?.error || e.message),
            );
        } finally {
            uploading = false;
            // Reset input
            target.value = "";
        }
    }

    async function deleteAddon(filename: string) {
        if (!confirm(`Voulez-vous vraiment supprimer ${filename} ?`)) {
            return;
        }

        try {
            await api.delete(
                `/api/servers/${serverId}/addons/${activeTab}/${filename}`,
            );
            await loadAddons();
            // Optional: Simple toast or just refresh is usually enough, but alert confirms action
        } catch (e: any) {
            console.error("Delete failed", e);
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

    function formatDate(dateStr: string): string {
        return new Date(dateStr).toLocaleString();
    }

    // Reload when tab changes
    $: if (activeTab) {
        loadAddons();
    }
</script>

<div class="space-y-6">
    <!-- Header / Tabs -->
    <div
        class="flex flex-col md:flex-row justify-between items-center gap-4 border-b border-gray-800 pb-4"
    >
        <div class="flex bg-gray-800/50 p-1 rounded-lg">
            {#each tabs as tab}
                <button
                    on:click={() => (activeTab = tab.id)}
                    class="flex items-center gap-2 px-4 py-2 rounded-md transition-all text-sm font-medium
                    {activeTab === tab.id
                        ? 'bg-blue-600 text-white shadow-lg'
                        : 'text-gray-400 hover:text-gray-200 hover:bg-gray-700/50'}"
                >
                    <svelte:component this={tab.icon} size={16} />
                    {tab.label}
                </button>
            {/each}
        </div>

        <button
            on:click={() => fileInput.click()}
            disabled={uploading}
            class="flex items-center gap-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-white px-4 py-2 rounded-lg font-medium transition-colors shadow-lg shadow-green-900/20"
        >
            {#if uploading}
                <Loader2 size={18} class="animate-spin" />
                Upload en cours...
            {:else}
                <Upload size={18} />
                Ajouter des {tabs.find((t) => t.id === activeTab)?.label}
            {/if}
        </button>
        <input
            bind:this={fileInput}
            type="file"
            multiple
            accept={tabs.find((t) => t.id === activeTab)?.accept}
            class="hidden"
            on:change={handleUpload}
        />
    </div>

    <!-- Content -->
    <div
        class="bg-slate-900/50 rounded-xl border border-gray-800 overflow-hidden min-h-[400px]"
    >
        {#if loading && !uploading}
            <div
                class="flex flex-col items-center justify-center h-64 text-gray-400 gap-3"
            >
                <Loader2 size={32} class="animate-spin text-blue-500" />
                <span>Chargement des fichiers...</span>
            </div>
        {:else if error}
            <div class="flex items-center justify-center h-64 text-red-400">
                {error}
            </div>
        {:else if files.length === 0}
            <div
                class="flex flex-col items-center justify-center h-64 text-gray-500 gap-4"
            >
                <div class="bg-gray-800/50 p-4 rounded-full">
                    <svelte:component
                        this={tabs.find((t) => t.id === activeTab)?.icon}
                        size={48}
                        class="opacity-50"
                    />
                </div>
                <p>
                    Aucun {tabs
                        .find((t) => t.id === activeTab)
                        ?.label.toLowerCase()
                        .slice(0, -1)} installé.
                </p>
            </div>
        {:else}
            <div class="overflow-x-auto min-h-[300px]">
                <table class="w-full text-left border-collapse">
                    <thead>
                        <tr
                            class="bg-gray-800/30 text-gray-400 text-xs uppercase border-b border-gray-800"
                        >
                            <th class="p-4 font-medium">Nom</th>
                            <th class="p-4 font-medium">Taille</th>
                            <th class="p-4 font-medium">Date de modification</th
                            >
                            <th class="p-4 font-medium text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-800">
                        {#each files as file}
                            <tr
                                class="hover:bg-gray-800/20 transition-colors group"
                            >
                                <td
                                    class="p-4 text-gray-200 font-medium flex items-center gap-3"
                                >
                                    <FileIcon size={18} class="text-blue-400" />
                                    {file.name}
                                </td>
                                <td class="p-4 text-gray-400 font-mono text-sm">
                                    {formatBytes(file.size)}
                                </td>
                                <td class="p-4 text-gray-400 text-sm">
                                    {formatDate(file.mod_time)}
                                </td>
                                <td class="p-4 text-right">
                                    <button
                                        on:click={() => deleteAddon(file.name)}
                                        class="p-2 text-gray-400 hover:text-red-400 hover:bg-red-900/20 rounded-lg transition-colors"
                                        title="Supprimer"
                                    >
                                        <Trash2 size={18} />
                                    </button>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>
