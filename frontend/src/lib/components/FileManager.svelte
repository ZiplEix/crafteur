<script lang="ts">
    import {
        Folder,
        File,
        Trash2,
        Upload,
        FolderPlus,
        ArrowUp,
        Archive,
        HardDrive,
    } from "lucide-svelte";
    import { api } from "$lib/api";
    import { onMount } from "svelte";

    let { serverId }: { serverId: string } = $props();

    interface FileInfo {
        name: string;
        is_dir: boolean;
        size: number;
        mod_time: string;
    }

    let currentPath = $state("/");
    let files: FileInfo[] = $state([]);
    let isUploading = $state(false);
    let loading = $state(true);
    let error: string | null = $state(null);
    let fileInput: HTMLInputElement;

    async function loadFiles(path: string) {
        loading = true;
        error = null;
        try {
            const res = await api.get(`/api/servers/${serverId}/files`, {
                params: { path },
            });
            files = res.data;
            // sort folders first
            files.sort((a, b) => {
                if (a.is_dir === b.is_dir) {
                    return a.name.localeCompare(b.name);
                }
                return a.is_dir ? -1 : 1;
            });
            currentPath = path;
        } catch (e) {
            console.error(e);
            error = "Impossible de charger les fichiers";
        } finally {
            loading = false;
        }
    }

    onMount(() => {
        loadFiles("/");
    });

    function navigate(folderName: string) {
        const newPath =
            currentPath === "/"
                ? `/${folderName}`
                : `${currentPath}/${folderName}`;
        loadFiles(newPath);
    }

    function navigateUp() {
        if (currentPath === "/") return;
        const parts = currentPath.split("/").filter((p) => p);
        parts.pop();
        const newPath = parts.length === 0 ? "/" : `/${parts.join("/")}`;
        loadFiles(newPath);
    }

    function navigateToBreadcrumb(index: number) {
        const parts = currentPath.split("/").filter((p) => p);
        const newPath =
            index === -1 ? "/" : `/${parts.slice(0, index + 1).join("/")}`;
        loadFiles(newPath);
    }

    async function handleCreateDir() {
        const name = prompt("Nom du nouveau dossier :");
        if (!name) return;
        try {
            await api.post(`/api/servers/${serverId}/files/directory`, {
                path: currentPath,
                name,
            });
            loadFiles(currentPath);
        } catch (e) {
            alert("Erreur lors de la création du dossier");
        }
    }

    async function handleDelete(filename: string) {
        if (!confirm(`Voulez-vous vraiment supprimer ${filename} ?`)) return;
        try {
            const targetPath =
                currentPath === "/"
                    ? `/${filename}`
                    : `${currentPath}/${filename}`;
            await api.delete(`/api/servers/${serverId}/files`, {
                params: { path: targetPath },
            });
            loadFiles(currentPath);
        } catch (e) {
            alert("Erreur lors de la suppression");
        }
    }

    async function handleUpload(e: Event) {
        const target = e.target as HTMLInputElement;
        if (!target.files || target.files.length === 0) return;
        const file = target.files[0];

        isUploading = true;
        const formData = new FormData();
        formData.append("file", file);
        formData.append("path", currentPath);

        try {
            await api.post(`/api/servers/${serverId}/files/upload`, formData, {
                headers: {
                    "Content-Type": "multipart/form-data",
                },
            });
            loadFiles(currentPath);
        } catch (e) {
            alert("Erreur lors de l'upload");
        } finally {
            isUploading = false;
            target.value = ""; // Reset input
        }
    }

    async function handleUnzip(filename: string) {
        if (!confirm(`Voulez-vous dézipper ${filename} ici ?`)) return;
        try {
            await api.post(`/api/servers/${serverId}/files/unzip`, {
                path: currentPath,
                filename,
            });
            loadFiles(currentPath);
            alert("Décompression terminée");
        } catch (e) {
            console.error(e);
            alert("Erreur lors de la décompression");
        }
    }

    function formatSize(bytes: number) {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }

    function formatDate(dateStr: string) {
        return new Date(dateStr).toLocaleString();
    }
</script>

<div class="flex flex-col h-full">
    <!-- Toolbar -->
    <div
        class="flex justify-between items-center mb-4 bg-gray-900 p-2 rounded-lg border border-gray-800"
    >
        <!-- Breadcrumbs -->
        <div
            class="flex items-center text-sm text-gray-300 overflow-x-auto whitespace-nowrap scrollbar-hide px-2"
        >
            <button
                onclick={() => navigateToBreadcrumb(-1)}
                class="hover:text-white hover:underline flex items-center cursor-pointer"
            >
                <HardDrive size={16} class="mr-1" />
                root
            </button>
            {#each currentPath.split("/").filter((p) => p) as part, index}
                <span class="mx-2 text-gray-600">/</span>
                <button
                    onclick={() => navigateToBreadcrumb(index)}
                    class="hover:text-white hover:underline cursor-pointer"
                >
                    {part}
                </button>
            {/each}
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-2">
            <input
                type="file"
                bind:this={fileInput}
                onchange={handleUpload}
                hidden
            />
            <button
                onclick={() => fileInput.click()}
                disabled={isUploading}
                class="flex items-center gap-2 bg-blue-600 hover:bg-blue-500 text-white px-3 py-1.5 rounded text-sm transition-colors cursor-pointer disabled:opacity-50"
            >
                {#if isUploading}
                    Scanning...
                {:else}
                    <Upload size={16} />
                    Upload
                {/if}
            </button>
            <button
                onclick={handleCreateDir}
                class="flex items-center gap-2 bg-gray-700 hover:bg-gray-600 text-white px-3 py-1.5 rounded text-sm transition-colors cursor-pointer"
            >
                <FolderPlus size={16} />
                Nouveau
            </button>
        </div>
    </div>

    <!-- File List -->
    <div
        class="border border-gray-800 rounded-lg bg-gray-900
        [&::-webkit-scrollbar]:w-2
        [&::-webkit-scrollbar-track]:bg-gray-900
        [&::-webkit-scrollbar-thumb]:bg-gray-700
        hover:[&::-webkit-scrollbar-thumb]:bg-gray-600"
    >
        {#if loading}
            <div class="p-8 text-center text-gray-500">Chargement...</div>
        {:else if error}
            <div class="p-8 text-center text-red-400">{error}</div>
        {:else}
            <table class="w-full text-sm text-left text-gray-300">
                <thead
                    class="text-xs text-gray-400 uppercase bg-gray-950 sticky top-0 z-10"
                >
                    <tr>
                        <th scope="col" class="px-6 py-3 w-10">Type</th>
                        <th scope="col" class="px-6 py-3">Nom</th>
                        <th scope="col" class="px-6 py-3">Taille</th>
                        <th scope="col" class="px-6 py-3">Date</th>
                        <th scope="col" class="px-6 py-3 text-right">Actions</th
                        >
                    </tr>
                </thead>
                <tbody>
                    {#if currentPath !== "/"}
                        <tr
                            class="border-b border-gray-800 hover:bg-gray-800/50 transition-colors"
                        >
                            <td class="px-6 py-3 text-center">
                                <ArrowUp size={18} class="text-blue-400" />
                            </td>
                            <td class="px-6 py-3">
                                <button
                                    onclick={navigateUp}
                                    class="font-medium text-blue-400 hover:underline cursor-pointer"
                                >
                                    ..
                                </button>
                            </td>
                            <td colspan="3"></td>
                        </tr>
                    {/if}

                    {#each files as file}
                        <tr
                            class="border-b border-gray-800 hover:bg-gray-800/50 transition-colors group"
                        >
                            <td class="px-6 py-3 text-center">
                                {#if file.is_dir}
                                    <Folder size={18} class="text-yellow-500" />
                                {:else}
                                    <File size={18} class="text-gray-500" />
                                {/if}
                            </td>
                            <td class="px-6 py-3 font-medium text-white">
                                {#if file.is_dir}
                                    <button
                                        onclick={() => navigate(file.name)}
                                        class="hover:underline hover:text-yellow-400 cursor-pointer text-left"
                                    >
                                        {file.name}
                                    </button>
                                {:else}
                                    <span>{file.name}</span>
                                {/if}
                            </td>
                            <td class="px-6 py-3 text-gray-400">
                                {file.is_dir ? "-" : formatSize(file.size)}
                            </td>
                            <td class="px-6 py-3 text-gray-500">
                                {formatDate(file.mod_time)}
                            </td>
                            <td class="px-6 py-3 text-right">
                                <div
                                    class="flex items-center justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity"
                                >
                                    {#if !file.is_dir && file.name.endsWith(".zip")}
                                        <button
                                            onclick={() =>
                                                handleUnzip(file.name)}
                                            class="text-blue-400 hover:text-blue-300 p-1 rounded hover:bg-blue-400/10 cursor-pointer"
                                            title="Dézipper"
                                        >
                                            <Archive size={16} />
                                        </button>
                                    {/if}
                                    <button
                                        onclick={() => handleDelete(file.name)}
                                        class="text-red-400 hover:text-red-300 p-1 rounded hover:bg-red-400/10 cursor-pointer"
                                        title="Supprimer"
                                    >
                                        <Trash2 size={16} />
                                    </button>
                                </div>
                            </td>
                        </tr>
                    {/each}
                    {#if files.length === 0 && currentPath === "/"}
                        <tr>
                            <td
                                colspan="5"
                                class="px-6 py-8 text-center text-gray-500 italic"
                            >
                                Dossier vide
                            </td>
                        </tr>
                    {/if}
                </tbody>
            </table>
        {/if}
    </div>
</div>
