<script lang="ts">
    import { onMount, onDestroy, tick, type Component } from "svelte";
    import { page } from "$app/stores";
    import { api } from "$lib/api";
    import {
        Terminal,
        FileText,
        Calendar,
        Save,
        Folder,
        Settings,
        Users,
        Play,
        Square,
        RotateCw,
    } from "lucide-svelte";
    import { configGroups } from "$lib/config_schema";
    import FileManager from "$lib/components/FileManager.svelte";

    interface ServerDetail {
        id: string;
        name: string;
        type: string;
        port: number;
        status: string;
    }

    let server: ServerDetail | null = null;
    let loading: boolean = true;
    let error: string | null = null;

    let activeTab: string = "console";
    let logs: string[] = [];
    let commandInput: string = "";
    let properties: Record<string, string> = {};

    let ws: WebSocket | null = null;
    let consoleContainer: HTMLElement;

    const serverId = $page.params.id;

    const tabs = [
        { id: "console", label: "Console", icon: Terminal },
        { id: "log", label: "Log", icon: FileText },
        { id: "schedule", label: "Schedule", icon: Calendar },
        { id: "save", label: "Save", icon: Save },
        { id: "file", label: "File", icon: Folder },
        { id: "configuration", label: "Configuration", icon: Settings },
        { id: "player", label: "Player", icon: Users },
    ];

    async function fetchServer() {
        try {
            const res = await api.get(`/api/servers/${serverId}`);
            server = res.data;
        } catch (e: any) {
            error = e.message;
        } finally {
            loading = false;
        }
    }

    async function fetchProperties() {
        try {
            const res = await api.get(`/api/servers/${serverId}/properties`);
            properties = res.data;
        } catch (e) {
            console.error("Failed to fetch properties", e);
            alert("Erreur lors du chargement de la configuration");
        }
    }

    async function saveProperties() {
        try {
            await api.post(`/api/servers/${serverId}/properties`, properties);
            alert(
                "Configuration sauvegardée ! Redémarrez le serveur pour appliquer.",
            );
        } catch (e) {
            console.error("Failed to save properties", e);
            alert("Erreur lors de la sauvegarde");
        }
    }

    async function sendAction(action: string) {
        if (!server) return;
        try {
            await api.post(`/api/servers/${server.id}/${action}`);
            // Status will be updated via WebSocket
            if (action === "start") {
                server.status = "STARTING";
            } else if (action === "stop") {
                server.status = "STOPPING";
            }
        } catch (e: any) {
            console.error(`Failed to ${action} server`, e);
            alert(`Failed to ${action} server: ` + e.message);
        }
    }

    async function sendCommand() {
        if (!commandInput.trim()) return;
        const cmd = commandInput;
        commandInput = ""; // Clear early for better UX

        try {
            await api.post(
                `/api/servers/${serverId}/command`,
                { command: cmd },
                {
                    headers: {
                        "Content-Type": "application/x-www-form-urlencoded",
                    },
                },
            );
            // Ideally the command echo comes back via WS
        } catch (e: any) {
            console.error("Failed to send command", e);
            logs = [...logs, `Error sending command: ${e.message}`];
        }
    }

    onMount(() => {
        fetchServer();
        connectWS();
    });

    onDestroy(() => {
        if (ws) {
            ws.close();
        }
    });

    function connectWS() {
        if (ws) {
            ws.close();
        }

        const wsUrl = `ws://localhost:8080/api/servers/${serverId}/ws`;

        ws = new WebSocket(wsUrl);

        ws.onopen = () => {
            logs = [...logs, "--- Connected to Server Console ---"];
        };

        ws.onmessage = (event) => {
            try {
                const msg = JSON.parse(event.data);
                if (msg.type === "log") {
                    logs = [...logs, msg.data];
                    if (activeTab === "console") {
                        scrollToBottom();
                    }
                } else if (msg.type === "status" && server) {
                    server.status = msg.data;
                }
            } catch (e) {
                console.error("Failed to parse WS message", e);
            }
        };

        ws.onclose = () => {
            logs = [...logs, "--- Connection Closed ---"];
            // Reconnect logic could go here
        };

        ws.onerror = (e) => {
            console.error("WS Error", e);
            logs = [...logs, "--- Connection Error ---"];
        };
    }

    async function scrollToBottom() {
        await tick();
        if (consoleContainer) {
            consoleContainer.scrollTop = consoleContainer.scrollHeight;
        }
    }

    $: if (activeTab === "configuration") {
        fetchProperties();
    }
</script>

<div class="container mx-auto p-6 max-w-7xl">
    {#if loading}
        <div class="text-white">Chargement...</div>
    {:else if error}
        <div class="text-red-500">Erreur: {error}</div>
    {:else if server}
        <!-- Header -->
        <div
            class="flex flex-col md:flex-row justify-between items-start md:items-center mb-6 gap-4"
        >
            <div class="flex items-center gap-4">
                <h1 class="text-3xl font-bold text-white">{server.name}</h1>
                <span
                    class="px-3 py-1 rounded-full text-sm font-semibold
                    {server.status === 'RUNNING'
                        ? 'bg-green-500/20 text-green-400'
                        : server.status === 'STOPPED'
                          ? 'bg-red-500/20 text-red-400'
                          : 'bg-yellow-500/20 text-yellow-400'}"
                >
                    {server.status}
                </span>
                <span class="text-gray-400 text-sm">Port: {server.port}</span>
            </div>

            <div class="flex gap-2">
                {#if server.status === "STOPPED"}
                    <button
                        on:click={() => sendAction("start")}
                        class="flex items-center gap-2 bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-lg transition-colors font-medium cursor-pointer"
                    >
                        <Play size={18} />
                        Démarrer
                    </button>
                {:else}
                    <button
                        on:click={() => sendAction("stop")}
                        class="flex items-center gap-2 bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors font-medium cursor-pointer"
                    >
                        <Square size={18} />
                        Arrêter
                    </button>
                    {#if server.status === "RUNNING"}
                        <button
                            on:click={() =>
                                sendAction("stop").then(() =>
                                    setTimeout(() => sendAction("start"), 2000),
                                )}
                            class="flex items-center gap-2 bg-orange-600 hover:bg-orange-700 text-white px-4 py-2 rounded-lg transition-colors font-medium cursor-pointer"
                        >
                            <RotateCw size={18} />
                            Redémarrer
                        </button>
                    {/if}
                {/if}
            </div>
        </div>

        <!-- Navigation Tabs -->
        <div class="border-b border-gray-700 mb-6 overflow-x-auto">
            <div class="flex gap-2">
                {#each tabs as tab}
                    <button
                        on:click={() => (activeTab = tab.id)}
                        class="flex items-center gap-2 px-4 py-3 border-b-2 transition-colors whitespace-nowrap cursor-pointer
                        {activeTab === tab.id
                            ? 'border-green-500 text-white'
                            : 'border-transparent text-gray-400 hover:text-gray-200 hover:border-gray-600'}"
                    >
                        <svelte:component this={tab.icon} size={18} />
                        {tab.label}
                    </button>
                {/each}
            </div>
        </div>

        <!-- Content Area -->
        <div class="bg-gray-900 rounded-xl border border-gray-800 p-6">
            {#if activeTab === "console"}
                <div class="flex flex-col" style="height: 66vh;">
                    <div
                        bind:this={consoleContainer}
                        class="flex-1 bg-black rounded-t-lg p-4 overflow-y-auto font-mono text-sm space-y-1 border border-gray-800 border-b-0
                        [&::-webkit-scrollbar]:w-2
                        [&::-webkit-scrollbar-track]:bg-gray-900
                        [&::-webkit-scrollbar-thumb]:bg-gray-700
                        hover:[&::-webkit-scrollbar-thumb]:bg-gray-600"
                    >
                        {#if logs.length === 0}
                            <div class="text-gray-500 italic">
                                En attente de logs...
                            </div>
                        {/if}
                        {#each logs as log}
                            <div
                                class="wrap-break-word text-gray-300 hover:bg-gray-900/50 px-1 rounded"
                            >
                                {log}
                            </div>
                        {/each}
                    </div>

                    <form
                        on:submit|preventDefault={sendCommand}
                        class="flex gap-0"
                    >
                        <div
                            class="bg-black text-green-500 px-3 py-3 font-mono border-l border-b border-gray-800 rounded-bl-lg select-none"
                        >
                            &gt;
                        </div>
                        <input
                            type="text"
                            bind:value={commandInput}
                            placeholder="Envoyer une commande..."
                            class="flex-1 bg-black text-gray-200 p-3 font-mono border-b border-r border-gray-800 rounded-br-lg focus:outline-none focus:bg-gray-950 transition-colors"
                        />
                    </form>
                </div>
            {:else if activeTab === "file"}
                <div>
                    <FileManager serverId={server.id} />
                </div>
            {:else if activeTab === "log"}
                <div
                    class="flex flex-col items-center justify-center h-64 text-gray-400"
                >
                    <FileText size={48} class="mb-4 opacity-50" />
                    <p class="text-lg">
                        Consultation des fichiers de logs (Bientôt)
                    </p>
                </div>
            {:else if activeTab === "configuration"}
                <div class="space-y-8">
                    {#each configGroups as group}
                        <div>
                            <h3
                                class="text-xl font-semibold text-white mb-4 border-b border-gray-800 pb-2"
                            >
                                {group.title}
                            </h3>
                            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                                {#each group.fields as field}
                                    {@const f = field as any}
                                    <div class="flex flex-col gap-2">
                                        <label
                                            class="text-sm font-medium text-gray-400"
                                            for={f.key}
                                        >
                                            {f.label}
                                        </label>
                                        {#if f.type === "select"}
                                            <select
                                                id={f.key}
                                                bind:value={properties[f.key]}
                                                class="bg-slate-800 border border-gray-700 rounded-lg p-2.5 text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                                            >
                                                {#each f.options || [] as option}
                                                    <option value={option}
                                                        >{option}</option
                                                    >
                                                {/each}
                                            </select>
                                        {:else if f.type === "boolean"}
                                            <div
                                                class="flex items-center h-[42px]"
                                            >
                                                <button
                                                    on:click={() =>
                                                        (properties[f.key] =
                                                            properties[
                                                                f.key
                                                            ] === "true"
                                                                ? "false"
                                                                : "true")}
                                                    class="relative w-11 h-6 rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 focus:ring-offset-gray-900 cursor-pointer {properties[
                                                        f.key
                                                    ] === 'true'
                                                        ? 'bg-green-600'
                                                        : 'bg-gray-700'}"
                                                >
                                                    <span
                                                        class="absolute top-[2px] start-[2px] bg-white border-gray-300 border rounded-full h-5 w-5 transition-transform {properties[
                                                            f.key
                                                        ] === 'true'
                                                            ? 'translate-x-[20px] border-white'
                                                            : 'translate-x-0'}"
                                                    ></span>
                                                </button>
                                                <span
                                                    class="ms-3 text-sm font-medium {properties[
                                                        f.key
                                                    ] === 'true'
                                                        ? 'text-green-400'
                                                        : 'text-gray-400'}"
                                                >
                                                    {properties[f.key] ===
                                                    "true"
                                                        ? "Activé"
                                                        : "Désactivé"}
                                                </span>
                                            </div>
                                        {:else}
                                            <input
                                                type={f.type}
                                                id={f.key}
                                                bind:value={properties[f.key]}
                                                min={f.min}
                                                max={f.max}
                                                class="bg-slate-800 border border-gray-700 rounded-lg p-2.5 text-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500 outline-none transition-all"
                                            />
                                        {/if}
                                    </div>
                                {/each}
                            </div>
                        </div>
                    {/each}

                    <div class="sticky bottom-6 flex justify-end pt-4">
                        <button
                            on:click={saveProperties}
                            class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-xl shadow-lg transition-transform hover:scale-105 font-medium"
                        >
                            <Save size={20} />
                            Sauvegarder la configuration
                        </button>
                    </div>
                </div>
            {:else}
                <div
                    class="flex flex-col items-center justify-center h-64 text-gray-400"
                >
                    <Folder size={48} class="mb-4 opacity-50" />
                    <p class="text-lg">
                        Section {tabs.find((t) => t.id === activeTab)?.label} en
                        cours de développement
                    </p>
                </div>
            {/if}
        </div>
    {/if}
</div>
