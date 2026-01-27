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
        Search,
        HardDrive,
    } from "lucide-svelte";
    import type { AddonFile, AddonType } from "$lib/types/addons";
    import ModrinthBrowser from "./ModrinthBrowser.svelte";

    export let serverId: string;
    export let serverType: string;

    let activeTab: AddonType = "mods";
    // 'installed' or 'catalog'
    let viewMode: "installed" | "catalog" = "installed";

    let files: AddonFile[] = [];
    let loading = false;
    let uploading = false;
    let error: string | null = null;
    let fileInput: HTMLInputElement;

    $: availableTabs = (() => {
        const t: { id: AddonType; label: string; icon: any; accept: string }[] =
            [];

        // Logique Fabric (Mods)
        if (serverType === "fabric") {
            t.push({
                id: "mods",
                label: "Mods",
                icon: Package,
                accept: ".jar",
            });
        }

        // Logique Paper (Plugins)
        if (serverType === "paper") {
            t.push({
                id: "plugins",
                label: "Plugins",
                icon: Zap,
                accept: ".jar",
            });
        }

        // Logique Universelle (Datapacks)
        t.push({
            id: "datapacks",
            label: "Datapacks",
            icon: Box,
            accept: ".zip",
        });

        return t;
    })();

    // Security: Auto-switch if active tab is invalid for current server type
    $: if (
        availableTabs.length > 0 &&
        !availableTabs.find((t) => t.id === activeTab)
    ) {
        activeTab = availableTabs[0].id;
    }

    async function loadAddons() {
        if (viewMode === "catalog") return; // Don't load files if in catalog mode

        loading = true;
        error = null;
        try {
            const res = await api.get(
                `/api/servers/${serverId}/addons/${activeTab}`,
            );
            files = res.data;
        } catch (e: any) {
            console.error("Failed to load addons", e);
            error = "Failed to load addons.";
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

            alert(`${formFiles.length} files added successfully!`);
            await loadAddons();
        } catch (e: any) {
            console.error("Upload failed", e);
            alert("Upload failed: " + (e.response?.data?.error || e.message));
        } finally {
            uploading = false;
            // Reset input
            target.value = "";
        }
    }

    async function deleteAddon(filename: string) {
        if (!confirm(`Are you sure you want to delete ${filename}?`)) {
            return;
        }

        try {
            await api.delete(
                `/api/servers/${serverId}/addons/${activeTab}/${filename}`,
            );
            await loadAddons();
        } catch (e: any) {
            console.error("Delete failed", e);
            alert("Delete failed: " + (e.response?.data?.error || e.message));
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

    function switchView(mode: "installed" | "catalog") {
        viewMode = mode;
        if (mode === "installed") {
            loadAddons();
        }
    }

    // Reload when tab changes
    $: if (activeTab) {
        if (
            activeTab !== "mods" &&
            activeTab !== "plugins" &&
            activeTab !== "datapacks"
        ) {
            viewMode = "installed";
        }
        loadAddons();
    }
</script>

<div class="space-y-6">
    <!-- Header / Tabs -->
    <div
        class="flex flex-col md:flex-row justify-between items-center gap-4 border-b border-gray-800 pb-4"
    >
        <div class="flex items-center gap-4">
            <div class="flex bg-gray-800/50 p-1 rounded-lg">
                {#each availableTabs as tab}
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

            {#if activeTab === "mods" || activeTab === "plugins" || activeTab === "datapacks"}
                <div
                    class="flex bg-gray-800/50 p-1 rounded-lg border-l border-gray-700 pl-1 ml-2"
                >
                    <button
                        on:click={() => switchView("installed")}
                        class="flex items-center gap-2 px-3 py-1.5 rounded-md transition-all text-xs font-medium
                        {viewMode === 'installed'
                            ? 'bg-gray-700 text-white'
                            : 'text-gray-400 hover:text-gray-200'}"
                    >
                        <HardDrive size={14} />
                        Installed
                    </button>
                    <button
                        on:click={() => switchView("catalog")}
                        class="flex items-center gap-2 px-3 py-1.5 rounded-md transition-all text-xs font-medium
                        {viewMode === 'catalog'
                            ? 'bg-emerald-600 text-white shadow-lg'
                            : 'text-gray-400 hover:text-gray-200'}"
                    >
                        <Search size={14} />
                        Catalog
                    </button>
                </div>
            {/if}
        </div>

        {#if viewMode === "installed"}
            <div class="flex items-center gap-2">
                <button
                    on:click={() => fileInput.click()}
                    disabled={uploading}
                    class="flex items-center gap-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed text-white px-4 py-2 rounded-lg font-medium transition-colors shadow-lg shadow-green-900/20"
                >
                    {#if uploading}
                        <Loader2 size={18} class="animate-spin" />
                        Uploading...
                    {:else}
                        <Upload size={18} />
                        Add files
                    {/if}
                </button>
                <input
                    bind:this={fileInput}
                    type="file"
                    multiple
                    accept={availableTabs.find((t) => t.id === activeTab)
                        ?.accept}
                    class="hidden"
                    on:change={handleUpload}
                />
            </div>
        {/if}
    </div>

    <!-- Content -->
    {#if viewMode === "installed"}
        <div
            class="bg-slate-900/50 rounded-xl border border-gray-800 overflow-hidden min-h-[400px]"
        >
            {#if loading && !uploading}
                <div
                    class="flex flex-col items-center justify-center h-64 text-gray-400 gap-3"
                >
                    <Loader2 size={32} class="animate-spin text-blue-500" />
                    <span>Loading files...</span>
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
                            this={availableTabs.find((t) => t.id === activeTab)
                                ?.icon}
                            size={48}
                            class="opacity-50"
                        />
                    </div>
                    <p>
                        No {availableTabs
                            .find((t) => t.id === activeTab)
                            ?.label.toLowerCase()
                            .slice(0, -1)} installed.
                    </p>
                </div>
            {:else}
                <div class="overflow-x-auto min-h-[300px]">
                    <table class="w-full text-left border-collapse">
                        <thead>
                            <tr
                                class="bg-gray-800/30 text-gray-400 text-xs uppercase border-b border-gray-800"
                            >
                                <th class="p-4 font-medium">Name</th>
                                <th class="p-4 font-medium">Size</th>
                                <th class="p-4 font-medium">Date Modified</th>
                                <th class="p-4 font-medium text-right"
                                    >Actions</th
                                >
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
                                        <FileIcon
                                            size={18}
                                            class="text-blue-400"
                                        />
                                        {file.name}
                                    </td>
                                    <td
                                        class="p-4 text-gray-400 font-mono text-sm"
                                    >
                                        {formatBytes(file.size)}
                                    </td>
                                    <td class="p-4 text-gray-400 text-sm">
                                        {formatDate(file.mod_time)}
                                    </td>
                                    <td class="p-4 text-right">
                                        <button
                                            on:click={() =>
                                                deleteAddon(file.name)}
                                            class="p-2 text-gray-400 hover:text-red-400 hover:bg-red-900/20 rounded-lg transition-colors"
                                            title="Delete"
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
    {:else if viewMode === "catalog" && (activeTab === "mods" || activeTab === "plugins" || activeTab === "datapacks")}
        <ModrinthBrowser
            {serverId}
            installedAddons={files}
            searchType={activeTab === "plugins"
                ? "plugin"
                : activeTab === "datapacks"
                  ? "datapack"
                  : "mod"}
        />
    {/if}
</div>
