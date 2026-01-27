<script lang="ts">
    import { onMount } from "svelte";
    import axios from "axios";
    import {
        Shield,
        Ban,
        Undo,
        ShieldCheck,
        ShieldX,
        DoorOpen,
        X,
    } from "lucide-svelte";
    import type { PlayerCache, OpEntry, BanEntry } from "$lib/types/players";

    let { serverId }: { serverId: string } = $props();

    let cache = $state<PlayerCache[]>([]);
    let ops = $state<OpEntry[]>([]);
    let bans = $state<BanEntry[]>([]);

    let isLoading = $state(true);
    let error = $state<string | null>(null);
    let successMessage = $state<string | null>(null);

    // Modal State
    let showReasonModal = $state(false);
    let targetPlayer = $state<string | null>(null);
    let targetAction = $state<"ban" | "kick" | null>(null);
    let reason = $state("");

    async function loadData() {
        isLoading = true;
        error = null;
        try {
            const [cacheRes, opsRes, bansRes] = await Promise.all([
                axios.get(
                    `http://localhost:8080/api/servers/${serverId}/players/cache`,
                    { withCredentials: true },
                ),
                axios.get(
                    `http://localhost:8080/api/servers/${serverId}/players/ops`,
                    { withCredentials: true },
                ),
                axios.get(
                    `http://localhost:8080/api/servers/${serverId}/players/banned`,
                    { withCredentials: true },
                ),
            ]);

            cache = cacheRes.data.sort((a: PlayerCache, b: PlayerCache) => {
                if (a.online && !b.online) return -1;
                if (!a.online && b.online) return 1;
                return (
                    new Date(b.expiresOn).getTime() -
                    new Date(a.expiresOn).getTime()
                );
            });
            ops = opsRes.data;
            bans = bansRes.data;
        } catch (e) {
            console.error("Failed to load player data", e);
            error = "Failed to load player data.";
        } finally {
            isLoading = false;
        }
    }

    async function handleAction(
        player: string,
        action: "op" | "deop" | "ban" | "pardon" | "kick",
        actionReason?: string,
    ) {
        try {
            await axios.post(
                `http://localhost:8080/api/servers/${serverId}/players/action`,
                {
                    player,
                    action,
                    reason: actionReason,
                },
                { withCredentials: true },
            );

            // Show success toast
            successMessage = `Action '${action}' performed successfully for ${player}.`;
            setTimeout(() => (successMessage = null), 3000);

            // Refetch after 1 second
            setTimeout(() => {
                loadData();
            }, 1000);
        } catch (e: any) {
            console.error(`Failed to execute ${action}`, e);
            error =
                e.response?.data?.error || `Error performing action ${action}`;
            setTimeout(() => (error = null), 5000);
        }
    }

    function openSanctionModal(player: string, action: "ban" | "kick") {
        targetPlayer = player;
        targetAction = action;
        reason = "";
        showReasonModal = true;
    }

    function closeSanctionModal() {
        showReasonModal = false;
        targetPlayer = null;
        targetAction = null;
        reason = "";
    }

    function confirmSanction() {
        if (targetPlayer && targetAction) {
            handleAction(targetPlayer, targetAction, reason);
            closeSanctionModal();
        }
    }

    onMount(() => {
        loadData();
    });
</script>

<div class="flex flex-col gap-6 p-4 relative">
    {#if successMessage}
        <div
            class="fixed top-4 right-4 bg-green-600 text-white px-4 py-2 rounded shadow-lg z-50 animate-bounce"
        >
            {successMessage}
        </div>
    {/if}

    {#if error}
        <div
            class="fixed top-4 right-4 bg-red-600 text-white px-4 py-2 rounded shadow-lg z-50"
        >
            {error}
        </div>
    {/if}

    {#if showReasonModal}
        <div
            class="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4"
        >
            <div
                class="bg-slate-800 border border-slate-700 rounded-lg p-6 w-full max-w-md shadow-2xl"
            >
                <div class="flex justify-between items-center mb-4">
                    <h3 class="text-xl font-bold text-white">
                        Punish {targetPlayer}
                    </h3>
                    <button
                        onclick={closeSanctionModal}
                        class="text-gray-400 hover:text-white cursor-pointer"
                        aria-label="Close"
                    >
                        <X size={24} />
                    </button>
                </div>
                <div class="mb-6">
                    <label
                        for="reason"
                        class="block text-sm font-medium text-slate-300 mb-2"
                        >Reason (Optional)</label
                    >
                    <input
                        type="text"
                        id="reason"
                        bind:value={reason}
                        class="w-full bg-slate-900 border border-slate-700 rounded p-2 text-white focus:outline-none focus:border-blue-500"
                        placeholder="Spam, Grief, etc."
                    />
                </div>
                <div class="flex justify-end gap-3">
                    <button
                        onclick={closeSanctionModal}
                        class="px-4 py-2 rounded bg-slate-700 text-white hover:bg-slate-600 transition-colors cursor-pointer"
                    >
                        Cancel
                    </button>
                    <button
                        onclick={confirmSanction}
                        class="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-700 transition-colors font-medium cursor-pointer"
                    >
                        Confirm
                    </button>
                </div>
            </div>
        </div>
    {/if}

    <div class="flex flex-col gap-8">
        <!-- Section 1: History (Cache) -->
        <div
            class="bg-slate-800 rounded-lg p-4 flex flex-col shadow-lg border border-slate-700"
        >
            <h2
                class="text-xl font-bold mb-4 text-slate-200 border-b border-slate-700 pb-2"
            >
                History
            </h2>
            <div class="overflow-y-auto flex-1 pr-2 space-y-2">
                {#if cache.length === 0}
                    <p class="text-gray-500 italic text-center mt-4">
                        No players in cache.
                    </p>
                {/if}
                {#each cache as entry}
                    <div
                        class="flex items-center justify-between p-3 bg-slate-900/50 rounded border border-slate-700 hover:border-slate-500 transition-colors {entry.online
                            ? 'border-l-4 border-l-green-500'
                            : ''}"
                    >
                        <div>
                            <div class="flex items-center gap-2">
                                {#if entry.online}
                                    <div
                                        class="w-2.5 h-2.5 rounded-full bg-green-500 shadow-[0_0_8px_#22c55e]"
                                        title="Online"
                                    ></div>
                                {/if}
                                <p class="font-medium text-slate-200">
                                    {entry.name}
                                </p>
                            </div>
                            <p class="text-xs text-slate-500">
                                Exp: {new Date(
                                    entry.expiresOn,
                                ).toLocaleDateString()}
                            </p>
                        </div>
                        <div class="flex gap-2">
                            <button
                                onclick={() => handleAction(entry.name, "op")}
                                class="p-1.5 text-green-400 hover:bg-green-400/10 rounded transition-colors cursor-pointer"
                                title="Promote to OP"
                            >
                                <Shield size={18} />
                            </button>
                            <button
                                onclick={() =>
                                    openSanctionModal(entry.name, "kick")}
                                class="p-1.5 text-orange-400 hover:bg-orange-400/10 rounded transition-colors cursor-pointer"
                                title="Kick"
                            >
                                <DoorOpen size={18} />
                            </button>
                            <button
                                onclick={() =>
                                    openSanctionModal(entry.name, "ban")}
                                class="p-1.5 text-red-400 hover:bg-red-400/10 rounded transition-colors cursor-pointer"
                                title="Ban"
                            >
                                <Ban size={18} />
                            </button>
                        </div>
                    </div>
                {/each}
            </div>
        </div>

        <!-- Section 2: Operators -->
        <div
            class="bg-slate-800 rounded-lg p-4 flex flex-col shadow-lg border border-slate-700"
        >
            <h2
                class="text-xl font-bold mb-4 text-slate-200 border-b border-slate-700 pb-2"
            >
                Operators
            </h2>
            <div class="overflow-y-auto flex-1 pr-2 space-y-2">
                {#if ops.length === 0}
                    <p class="text-gray-500 italic text-center mt-4">
                        No operators.
                    </p>
                {/if}
                {#each ops as entry}
                    <div
                        class="flex items-center justify-between p-3 bg-slate-900/50 rounded border border-slate-700 hover:border-slate-500 transition-colors"
                    >
                        <div>
                            <div class="flex items-center gap-2">
                                <ShieldCheck size={16} class="text-green-500" />
                                <p class="font-medium text-slate-200">
                                    {entry.name}
                                </p>
                            </div>
                            <p class="text-xs text-slate-500">
                                Level: {entry.level}
                            </p>
                        </div>
                        <button
                            onclick={() => handleAction(entry.name, "deop")}
                            class="p-1.5 text-orange-400 hover:bg-orange-400/10 rounded transition-colors cursor-pointer"
                            title="Deop"
                        >
                            <ShieldX size={18} />
                        </button>
                    </div>
                {/each}
            </div>
        </div>

        <!-- Section 3: Banned -->
        <div
            class="bg-slate-800 rounded-lg p-4 flex flex-col shadow-lg border border-slate-700"
        >
            <h2
                class="text-xl font-bold mb-4 text-slate-200 border-b border-slate-700 pb-2"
            >
                Banned
            </h2>
            <div class="overflow-y-auto flex-1 pr-2 space-y-2">
                {#if bans.length === 0}
                    <p class="text-gray-500 italic text-center mt-4">
                        No banned players.
                    </p>
                {/if}
                {#each bans as entry}
                    <div
                        class="flex items-center justify-between p-3 bg-slate-900/50 rounded border border-slate-700 hover:border-slate-500 transition-colors"
                    >
                        <div>
                            <p
                                class="font-medium text-red-300 line-through decoration-red-500/50"
                            >
                                {entry.name}
                            </p>
                            <p class="text-xs text-slate-400 italic">
                                "{entry.reason}"
                            </p>
                            <p class="text-xs text-slate-500">
                                By: {entry.source}
                            </p>
                        </div>
                        <button
                            onclick={() => handleAction(entry.name, "pardon")}
                            class="p-1.5 text-blue-400 hover:bg-blue-400/10 rounded transition-colors cursor-pointer"
                            title="Pardon"
                        >
                            <Undo size={18} />
                        </button>
                    </div>
                {/each}
            </div>
        </div>
    </div>
</div>
