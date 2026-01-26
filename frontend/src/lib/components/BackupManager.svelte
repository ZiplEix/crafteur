<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { Archive, Download, Trash2, Plus, HardDrive } from "lucide-svelte";
    import type { AxiosError } from "axios";

    // Receives prop serverId
    let { serverId }: { serverId: string } = $props();

    interface BackupEntry {
        name: string;
        size: number;
        created_at: string;
    }

    let backups: BackupEntry[] = $state([]);
    let loading: boolean = $state(true);
    let error: string | null = $state(null);
    let isCreating: boolean = $state(false);

    async function fetchBackups() {
        loading = true;
        error = null;
        try {
            const res = await api.get(`/api/servers/${serverId}/backups`);
            backups = res.data;
        } catch (e: any) {
            console.error("Failed to fetch backups", e);
            error = "Impossible de charger les sauvegardes";
        } finally {
            loading = false;
        }
    }

    async function createBackup() {
        if (isCreating) return;
        isCreating = true;
        error = null;
        try {
            await api.post(`/api/servers/${serverId}/backups`);
            alert("Sauvegarde créée avec succès !");
            await fetchBackups();
        } catch (e: any) {
            const ae = e as AxiosError;
            console.error("Failed to create backup", e);
            const data = ae.response?.data as any;
            alert(
                `Erreur lors de la création de la sauvegarde: ${data?.error || e.message}`,
            );
        } finally {
            isCreating = false;
        }
    }

    async function deleteBackup(filename: string) {
        if (
            !confirm(
                `Voulez-vous vraiment supprimer la sauvegarde ${filename} ?`,
            )
        )
            return;

        try {
            await api.delete(`/api/servers/${serverId}/backups/${filename}`);
            await fetchBackups();
        } catch (e: any) {
            const ae = e as AxiosError;
            console.error("Failed to delete backup", e);
            const data = ae.response?.data as any;
            alert(`Erreur lors de la suppression: ${data?.error || e.message}`);
        }
    }

    function formatBytes(bytes: number) {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }

    function formatDate(dateStr: string) {
        return new Date(dateStr).toLocaleString();
    }

    // Construct download link
    function getDownloadLink(filename: string) {
        // Assuming API base URL is relative or same origin, but we need full path if client is on different port/domain in dev
        // However, axios usually has base URL configured. For a simple link, we might need the full URL if frontend/backend are separate.
        // Given setup, backend is :8080.
        return `http://localhost:8080/api/servers/${serverId}/backups/${filename}`;
    }

    onMount(() => {
        fetchBackups();
    });
</script>

<div class="space-y-6">
    <!-- Header / Actions -->
    <div
        class="flex justify-between items-center bg-gray-900 p-4 rounded-lg border border-gray-800"
    >
        <div class="flex items-center gap-3">
            <div class="bg-blue-500/20 p-2 rounded-lg text-blue-400">
                <Archive size={24} />
            </div>
            <div>
                <h3 class="font-semibold text-white">Sauvegardes</h3>
                <p class="text-xs text-gray-400">
                    Gérez les snapshots de votre serveur
                </p>
            </div>
        </div>

        <div class="flex items-center gap-4">
            <button
                onclick={createBackup}
                disabled={isCreating}
                class="flex items-center gap-2 bg-green-600 hover:bg-green-700 disabled:bg-green-600/50 disabled:cursor-not-allowed text-white px-4 py-2 rounded-lg transition-all font-medium shadow-lg hover:shadow-green-900/20 active:scale-95"
            >
                {#if isCreating}
                    <div
                        class="w-4 h-4 border-2 border-white/50 border-t-white rounded-full animate-spin"
                    ></div>
                    <span>Création en cours...</span>
                {:else}
                    <Plus size={18} />
                    <span>Créer une sauvegarde</span>
                {/if}
            </button>
        </div>
    </div>

    <!-- List -->
    <div
        class="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden shadow-sm"
    >
        {#if loading}
            <div class="p-12 text-center">
                <div
                    class="inline-block w-8 h-8 border-4 border-gray-700 border-t-blue-500 rounded-full animate-spin mb-4"
                ></div>
                <p class="text-gray-500">Chargement des sauvegardes...</p>
            </div>
        {:else if error}
            <div
                class="p-12 text-center text-red-400 flex flex-col items-center gap-2"
            >
                <Trash2 size={32} class="opacity-50" />
                <p>{error}</p>
                <button
                    onclick={fetchBackups}
                    class="text-blue-400 hover:underline text-sm mt-2"
                    >Réessayer</button
                >
            </div>
        {:else if backups.length === 0}
            <div
                class="p-16 text-center text-gray-500 flex flex-col items-center gap-3"
            >
                <HardDrive size={48} class="opacity-20" />
                <p class="text-lg font-medium">Aucune sauvegarde</p>
                <p class="text-sm">
                    Créez votre première sauvegarde pour sécuriser vos données.
                </p>
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left text-sm text-gray-400">
                    <thead
                        class="bg-gray-950 text-gray-200 uppercase text-xs font-semibold"
                    >
                        <tr>
                            <th scope="col" class="px-6 py-4">Nom</th>
                            <th scope="col" class="px-6 py-4"
                                >Date de création</th
                            >
                            <th scope="col" class="px-6 py-4">Taille</th>
                            <th scope="col" class="px-6 py-4 text-right"
                                >Actions</th
                            >
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-800">
                        {#each backups as backup}
                            <tr
                                class="hover:bg-gray-800/50 transition-colors group"
                            >
                                <td
                                    class="px-6 py-4 font-medium text-white flex items-center gap-3"
                                >
                                    <Archive size={16} class="text-blue-500" />
                                    {backup.name}
                                </td>
                                <td class="px-6 py-4">
                                    {formatDate(backup.created_at)}
                                </td>
                                <td class="px-6 py-4 font-mono text-gray-500">
                                    {formatBytes(backup.size)}
                                </td>
                                <td class="px-6 py-4 text-right">
                                    <div
                                        class="flex items-center justify-end gap-2"
                                    >
                                        <a
                                            href={getDownloadLink(backup.name)}
                                            download
                                            class="p-2 text-blue-400 hover:text-white hover:bg-blue-500 rounded-lg transition-colors cursor-pointer"
                                            title="Télécharger"
                                        >
                                            <Download size={18} />
                                        </a>
                                        <button
                                            onclick={() =>
                                                deleteBackup(backup.name)}
                                            class="p-2 text-red-400 hover:text-white hover:bg-red-500 rounded-lg transition-colors cursor-pointer"
                                            title="Supprimer"
                                        >
                                            <Trash2 size={18} />
                                        </button>
                                    </div>
                                </td>
                            </tr>
                        {/each}
                    </tbody>
                </table>
            </div>
        {/if}
    </div>
</div>
