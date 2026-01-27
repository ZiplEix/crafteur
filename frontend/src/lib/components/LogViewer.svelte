<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { FileText, Search, RefreshCw, AlertCircle } from "lucide-svelte";

    let { serverId }: { serverId: string } = $props();

    interface LogFileEntry {
        name: string;
        size: number;
        mod_time: string;
    }

    let logs: LogFileEntry[] = $state([]);
    let loadingList: boolean = $state(true);
    let errorList: string | null = $state(null);

    let selectedFile: string | null = $state(null);
    let rawContent: string = $state("");
    let loadingContent: boolean = $state(false);
    let errorContent: string | null = $state(null);

    let searchQuery: string = $state("");
    let useRegex: boolean = $state(false);
    let ignoreCase: boolean = $state(true);
    let regexError: boolean = $state(false);

    // Derived filtered lines based on search
    let filteredLines = $derived.by(() => {
        if (!rawContent) return [];
        const lines = rawContent.split("\n");

        if (!searchQuery) return lines;

        try {
            let regex: RegExp;
            if (useRegex) {
                regex = new RegExp(searchQuery, ignoreCase ? "i" : "");
                regexError = false; // Reset error if successful
            } else {
                // Escape special regex chars for simple text search if we were to use regex implementation for both,
                // but here for simple search we can just use includes/toLowerCase
                const query = ignoreCase
                    ? searchQuery.toLowerCase()
                    : searchQuery;
                return lines.filter((line) => {
                    const l = ignoreCase ? line.toLowerCase() : line;
                    return l.includes(query);
                });
            }

            return lines.filter((line) => regex.test(line));
        } catch (e) {
            regexError = true;
            return lines; // Return all lines or empty? Let's return all but show error on input
        }
    });

    async function fetchLogs() {
        loadingList = true;
        errorList = null;
        try {
            const res = await api.get(`/api/servers/${serverId}/logs`);
            logs = res.data;
            // Ensure latest.log is first if present
            const latestIndex = logs.findIndex((l) => l.name === "latest.log");
            if (latestIndex > 0) {
                const latest = logs.splice(latestIndex, 1)[0];
                logs.unshift(latest);
            }
        } catch (e: any) {
            console.error("Failed to fetch logs", e);
            errorList = "Impossible de charger la liste des logs";
        } finally {
            loadingList = false;
        }
    }

    async function loadLogContent(filename: string) {
        if (selectedFile === filename && rawContent) return;

        selectedFile = filename;
        loadingContent = true;
        errorContent = null;
        rawContent = "";

        try {
            const res = await api.get(`/api/servers/${serverId}/logs/content`, {
                params: { filename },
                transformResponse: [(data) => data], // Force raw text
            });
            rawContent = res.data;
        } catch (e: any) {
            console.error("Failed to load log content", e);
            errorContent = "Impossible de lire le fichier";
        } finally {
            loadingContent = false;
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

    onMount(() => {
        fetchLogs();
    });
</script>

<div
    class="flex h-[calc(100vh-200px)] border border-gray-800 rounded-xl overflow-hidden bg-gray-900"
>
    <!-- Sidebar -->
    <div class="w-1/4 border-r border-gray-800 flex flex-col bg-gray-950">
        <div
            class="p-4 border-b border-gray-800 flex justify-between items-center"
        >
            <h3 class="font-semibold text-gray-200">Fichiers</h3>
            <button
                onclick={fetchLogs}
                class="text-gray-400 hover:text-white transition-colors cursor-pointer"
            >
                <RefreshCw size={16} />
            </button>
        </div>

        <div
            class="flex-1 overflow-y-auto
            [&::-webkit-scrollbar]:w-2
            [&::-webkit-scrollbar-track]:bg-gray-950
            [&::-webkit-scrollbar-thumb]:bg-gray-800
            hover:[&::-webkit-scrollbar-thumb]:bg-gray-700"
        >
            {#if loadingList}
                <div class="p-4 text-center text-gray-500 text-sm">
                    Chargement...
                </div>
            {:else if errorList}
                <div class="p-4 text-center text-red-400 text-sm">
                    {errorList}
                </div>
            {:else if logs.length === 0}
                <div class="p-4 text-center text-gray-500 text-sm">
                    Aucun log trouvé
                </div>
            {:else}
                {#each logs as log}
                    <button
                        onclick={() => loadLogContent(log.name)}
                        class="w-full text-left px-4 py-3 border-b border-gray-900 hover:bg-gray-900 transition-colors flex flex-col gap-1
                         {selectedFile === log.name
                            ? 'bg-gray-900 border-l-2 border-l-blue-500'
                            : 'border-l-2 border-l-transparent text-gray-400'} cursor-pointer"
                    >
                        <div
                            class="font-medium text-sm truncate {selectedFile ===
                            log.name
                                ? 'text-white'
                                : ''}"
                        >
                            {log.name}
                        </div>
                        <div class="flex justify-between text-xs text-gray-500">
                            <span>{formatSize(log.size)}</span>
                            <span>{formatDate(log.mod_time).split(",")[0]}</span
                            >
                        </div>
                    </button>
                {/each}
            {/if}
        </div>
    </div>

    <!-- Main Content -->
    <div class="w-3/4 flex flex-col bg-gray-900">
        <!-- Toolbar -->
        <div
            class="p-3 border-b border-gray-800 bg-gray-950 flex items-center gap-4"
        >
            <div class="relative flex-1">
                <Search
                    size={16}
                    class="absolute left-3 top-1/2 -translate-y-1/2 text-gray-500"
                />
                <input
                    type="text"
                    bind:value={searchQuery}
                    placeholder="Rechercher..."
                    class="w-full bg-gray-900 border {regexError
                        ? 'border-red-500'
                        : 'border-gray-800'} rounded-lg pl-9 pr-3 py-1.5 text-sm text-gray-200 focus:ring-1 focus:ring-blue-500 outline-none transition-all placeholder:text-gray-600"
                />
            </div>

            <div
                class="flex items-center gap-3 text-sm text-gray-300 select-none"
            >
                <label
                    class="flex items-center gap-2 cursor-pointer hover:text-white"
                >
                    <input
                        type="checkbox"
                        bind:checked={useRegex}
                        class="rounded bg-gray-800 border-gray-700 text-blue-600 focus:ring-offset-gray-900"
                    />
                    Regex
                </label>
                <label
                    class="flex items-center gap-2 cursor-pointer hover:text-white"
                >
                    <input
                        type="checkbox"
                        bind:checked={ignoreCase}
                        class="rounded bg-gray-800 border-gray-700 text-blue-600 focus:ring-offset-gray-900"
                    />
                    Ignorer la casse
                </label>
            </div>
        </div>

        <!-- Viewer -->
        <div
            class="flex-1 overflow-y-auto p-4 font-mono text-xs md:text-sm text-gray-300 whitespace-pre-wrap
            [&::-webkit-scrollbar]:w-2
            [&::-webkit-scrollbar-track]:bg-gray-900
            [&::-webkit-scrollbar-thumb]:bg-gray-700
            hover:[&::-webkit-scrollbar-thumb]:bg-gray-600"
        >
            {#if loadingContent}
                <div
                    class="flex flex-col items-center justify-center h-full text-gray-500 gap-2"
                >
                    <div
                        class="w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"
                    ></div>
                    <span>Chargement du contenu...</span>
                </div>
            {:else if errorContent}
                <div
                    class="flex flex-col items-center justify-center h-full text-red-400 gap-2 opacity-80"
                >
                    <AlertCircle size={32} />
                    <span>{errorContent}</span>
                </div>
            {:else if !selectedFile}
                <div
                    class="flex flex-col items-center justify-center h-full text-gray-600 gap-2"
                >
                    <FileText size={48} class="opacity-20" />
                    <span>Sélectionnez un fichier pour voir son contenu</span>
                </div>
            {:else if filteredLines.length === 0}
                <div
                    class="flex flex-col items-center justify-center h-full text-gray-500 italic"
                >
                    Aucun résultat pour la recherche
                </div>
            {:else}
                {#each filteredLines as line}
                    <!-- Simple syntax highlighting logic using classes -->
                    <div
                        class="leading-relaxed break-all hover:bg-gray-800/30 px-1 rounded
                        {line.includes('[ERROR]') ? 'text-red-400' : ''}
                        {line.includes('[WARN]') ? 'text-yellow-400' : ''}
                        {line.includes('[INFO]') ? 'text-blue-300' : ''}
                        {line.includes('[FATAL]')
                            ? 'text-red-500 font-bold'
                            : ''}
                        {line.includes('Exception') ? 'text-red-300' : ''}
                    "
                    >
                        {line}
                    </div>
                {/each}
            {/if}
        </div>

        <!-- Footer / Status Bar -->
        {#if selectedFile}
            <div
                class="bg-gray-950 border-t border-gray-800 px-4 py-1 flex justify-between text-xs text-gray-500 font-mono"
            >
                <span>{selectedFile}</span>
                <span
                    >{filteredLines.length} ligne{filteredLines.length > 1
                        ? "s"
                        : ""}</span
                >
            </div>
        {/if}
    </div>
</div>
