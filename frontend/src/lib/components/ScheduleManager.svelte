<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import {
        Calendar,
        Trash2,
        Plus,
        Clock,
        Play,
        RotateCw,
        Square,
        Terminal,
        X,
        Save,
    } from "lucide-svelte";
    import type { AxiosError } from "axios";
    import {
        generateCronFromRecurring,
        generateCronFromInterval,
    } from "$lib/utils/cronHelper";

    let { serverId }: { serverId: string } = $props();

    interface ScheduledTask {
        id: string;
        name: string;
        action: string;
        cron_expression: string;
        last_run: string;
        next_run: string;
        one_shot: boolean;
    }

    let tasks: ScheduledTask[] = $state([]);
    let loading: boolean = $state(true);
    let error: string | null = $state(null);

    // Modal State
    let isModalOpen = $state(false);
    let isCreating = $state(false);

    // Form State
    let name = $state("");
    let action = $state("start");
    let payload = $state("");
    let oneShot = $state(false);
    let triggerType = $state<"recur" | "interval" | "custom">("recur");

    // Trigger - Recurring
    let recurTime = $state("10:00");
    let recurDays: number[] = $state([]); // 0-6 Sun-Sat

    // Trigger - Interval
    let intervalVal = $state(1);
    let intervalUnit = $state<"m" | "h">("h");

    // Trigger - Custom
    let customCron = $state("");

    async function fetchTasks() {
        loading = true;
        try {
            const res = await api.get(`/api/servers/${serverId}/tasks`);
            tasks = res.data || [];
        } catch (e: any) {
            console.error("Failed to fetch tasks", e);
            error = "Failed to load tasks";
        } finally {
            loading = false;
        }
    }

    async function createTask() {
        isCreating = true;
        let cronExpr = "";

        if (triggerType === "recur") {
            // If no days selected, assume every day (*)
            cronExpr = generateCronFromRecurring(
                recurTime,
                recurDays.length > 0 ? recurDays : [0, 1, 2, 3, 4, 5, 6],
            );
        } else if (triggerType === "interval") {
            cronExpr = generateCronFromInterval(intervalVal, intervalUnit);
        } else {
            cronExpr = customCron;
        }

        try {
            await api.post(`/api/servers/${serverId}/tasks`, {
                name,
                action,
                payload,
                cron_expression: cronExpr,
                one_shot: oneShot,
            });
            isModalOpen = false;
            resetForm();
            await fetchTasks();
        } catch (e: any) {
            const ae = e as AxiosError;
            console.error("Failed to create task", e);
            const data = ae.response?.data as any;
            alert(`Failed to create task: ${data?.error || e.message}`);
        } finally {
            isCreating = false;
        }
    }

    async function deleteTask(id: string) {
        if (!confirm("Are you sure you want to delete this task?")) return;
        try {
            await api.delete(`/api/servers/${serverId}/tasks/${id}`);
            await fetchTasks();
        } catch (e: any) {
            const ae = e as AxiosError;
            const data = ae.response?.data as any;
            alert(`Failed to delete task: ${data?.error || e.message}`);
        }
    }

    function resetForm() {
        name = "";
        action = "start";
        payload = "";
        oneShot = false;
        triggerType = "recur";
        recurTime = "10:00";
        recurDays = [];
        intervalVal = 1;
        intervalUnit = "h";
        customCron = "";
    }

    function formatDate(dateStr: string) {
        if (dateStr.startsWith("0001")) return "Never";
        return new Date(dateStr).toLocaleString();
    }

    function toggleDay(day: number) {
        if (recurDays.includes(day)) {
            recurDays = recurDays.filter((d) => d !== day);
        } else {
            recurDays = [...recurDays, day];
        }
    }

    const daysOfWeek = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

    onMount(() => {
        fetchTasks();
    });
</script>

<div class="space-y-6">
    <!-- Header -->
    <div
        class="flex justify-between items-center bg-gray-900 p-4 rounded-lg border border-gray-800"
    >
        <div class="flex items-center gap-3">
            <div class="bg-purple-500/20 p-2 rounded-lg text-purple-400">
                <Calendar size={24} />
            </div>
            <div>
                <h3 class="font-semibold text-white">Scheduled Tasks</h3>
                <p class="text-xs text-gray-400">
                    Automate your server management
                </p>
            </div>
        </div>

        <button
            onclick={() => (isModalOpen = true)}
            class="flex items-center gap-2 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-all font-medium shadow-lg hover:shadow-blue-900/20 active:scale-95 cursor-pointer"
        >
            <Plus size={18} />
            <span>New Task</span>
        </button>
    </div>

    <!-- List -->
    <div
        class="bg-gray-900 border border-gray-800 rounded-xl overflow-hidden shadow-sm"
    >
        {#if loading}
            <div class="p-12 text-center text-gray-500">Loading tasks...</div>
        {:else if tasks.length === 0}
            <div
                class="p-16 text-center text-gray-500 flex flex-col items-center gap-3"
            >
                <Clock size={48} class="opacity-20" />
                <p class="text-lg font-medium">No scheduled tasks</p>
            </div>
        {:else}
            <div class="overflow-x-auto">
                <table class="w-full text-left text-sm text-gray-400">
                    <thead
                        class="bg-gray-950 text-gray-200 uppercase text-xs font-semibold"
                    >
                        <tr>
                            <th class="px-6 py-4">Name</th>
                            <th class="px-6 py-4">Action</th>
                            <th class="px-6 py-4">Frequency</th>
                            <th class="px-6 py-4">Last Run</th>
                            <th class="px-6 py-4">Next Run</th>
                            <th class="px-6 py-4 text-right">Actions</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-800">
                        {#each tasks as task}
                            <tr class="hover:bg-gray-800/50 transition-colors">
                                <td class="px-6 py-4 font-medium text-white"
                                    >{task.name}</td
                                >
                                <td class="px-6 py-4">
                                    <span
                                        class={`inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-xs font-medium border
                                        ${
                                            task.action === "start"
                                                ? "bg-green-500/10 text-green-400 border-green-500/20"
                                                : task.action === "stop"
                                                  ? "bg-red-500/10 text-red-400 border-red-500/20"
                                                  : task.action === "restart"
                                                    ? "bg-yellow-500/10 text-yellow-400 border-yellow-500/20"
                                                    : "bg-blue-500/10 text-blue-400 border-blue-500/20"
                                        }`}
                                    >
                                        {#if task.action === "start"}
                                            <Play size={10} />
                                        {:else if task.action === "stop"}
                                            <Square size={10} />
                                        {:else if task.action === "restart"}
                                            <RotateCw size={10} />
                                        {:else}
                                            <Terminal size={10} />
                                        {/if}
                                        {task.action.toUpperCase()}
                                    </span>
                                </td>
                                <td class="px-6 py-4 font-mono text-gray-400"
                                    >{task.cron_expression}</td
                                >
                                <td class="px-6 py-4"
                                    >{formatDate(task.last_run)}</td
                                >
                                <td class="px-6 py-4 text-blue-300"
                                    >{formatDate(task.next_run)}</td
                                >
                                <td class="px-6 py-4 text-right">
                                    <button
                                        onclick={() => deleteTask(task.id)}
                                        class="p-2 text-red-400 hover:text-white hover:bg-red-500 rounded-lg transition-colors cursor-pointer"
                                    >
                                        <Trash2 size={16} />
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

<!-- Modal -->
{#if isModalOpen}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm p-4"
    >
        <div
            class="bg-gray-900 border border-gray-800 rounded-xl w-full max-w-lg shadow-2xl flex flex-col max-h-[90vh]"
        >
            <div
                class="p-6 border-b border-gray-800 flex justify-between items-center"
            >
                <h3 class="text-xl font-semibold text-white">
                    New Scheduled Task
                </h3>
                <button
                    onclick={() => (isModalOpen = false)}
                    class="text-gray-400 hover:text-white cursor-pointer"
                >
                    <X size={20} />
                </button>
            </div>

            <div class="p-6 overflow-y-auto space-y-5">
                <!-- Name -->
                <div class="space-y-2">
                    <label class="text-sm font-medium text-gray-300"
                        >Task Name</label
                    >
                    <input
                        bind:value={name}
                        type="text"
                        placeholder="e.g., Daily Restart"
                        class="w-full bg-gray-950 border border-gray-800 rounded-lg px-4 py-2.5 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                    />
                </div>

                <!-- Action -->
                <div class="space-y-2">
                    <label class="text-sm font-medium text-gray-300"
                        >Action</label
                    >
                    <div class="grid grid-cols-4 gap-2">
                        {#each ["start", "restart", "stop", "command"] as act}
                            <button
                                onclick={() => (action = act)}
                                class={`px-3 py-2 rounded-lg text-sm font-medium border transition-all cursor-pointer ${action === act ? "bg-blue-600 border-blue-500 text-white" : "bg-gray-800 border-gray-700 text-gray-400 hover:bg-gray-700"}`}
                            >
                                {act.charAt(0).toUpperCase() + act.slice(1)}
                            </button>
                        {/each}
                    </div>
                </div>

                <!-- Payload if command -->
                {#if action === "command"}
                    <div class="space-y-2">
                        <label class="text-sm font-medium text-gray-300"
                            >Command(s)</label
                        >
                        <textarea
                            bind:value={payload}
                            rows={3}
                            placeholder="say Hello World&#10;save-all"
                            class="w-full bg-gray-950 border border-gray-800 rounded-lg px-4 py-2.5 font-mono text-sm text-white focus:ring-2 focus:ring-blue-500 outline-none"
                        ></textarea>
                        <p class="text-xs text-gray-500">
                            One command per line
                        </p>
                    </div>
                {/if}

                <!-- Trigger Type Tabs -->
                <div class="space-y-4 pt-2 border-t border-gray-800">
                    <label class="text-sm font-medium text-gray-300"
                        >Trigger Frequency</label
                    >
                    <div
                        class="flex bg-gray-950 p-1 rounded-lg border border-gray-800"
                    >
                        {#each ["recur", "interval", "custom"] as type}
                            <button
                                onclick={() => (triggerType = type as any)}
                                class={`flex-1 py-1.5 text-sm font-medium rounded-md transition-all cursor-pointer ${triggerType === type ? "bg-gray-800 text-white shadow-sm" : "text-gray-400 hover:text-gray-300"}`}
                            >
                                {type === "recur"
                                    ? "Recurring"
                                    : type === "interval"
                                      ? "Interval"
                                      : "Custom"}
                            </button>
                        {/each}
                    </div>

                    <div
                        class="bg-gray-950/50 p-4 rounded-lg border border-gray-800/50"
                    >
                        {#if triggerType === "recur"}
                            <div class="space-y-4">
                                <div class="flex items-center gap-4">
                                    <label class="text-sm text-gray-400"
                                        >At Time:</label
                                    >
                                    <input
                                        bind:value={recurTime}
                                        type="time"
                                        class="bg-gray-900 border border-gray-700 rounded px-2 py-1 text-white focus:ring-1 focus:ring-blue-500 outline-none"
                                    />
                                </div>
                                <div class="space-y-2">
                                    <label class="text-sm text-gray-400"
                                        >On Days:</label
                                    >
                                    <div class="flex flex-wrap gap-2">
                                        {#each daysOfWeek as day, i}
                                            <button
                                                onclick={() => toggleDay(i)}
                                                class={`w-10 h-10 rounded-full text-xs font-bold transition-all cursor-pointer ${recurDays.includes(i) ? "bg-blue-600 text-white shadow-lg shadow-blue-900/30 ring-2 ring-blue-400/50" : "bg-gray-800 text-gray-500 hover:bg-gray-700"}`}
                                            >
                                                {day}
                                            </button>
                                        {/each}
                                    </div>
                                    <p class="text-xs text-gray-500 mt-1">
                                        Select no days to run every day.
                                    </p>
                                </div>
                            </div>
                        {:else if triggerType === "interval"}
                            <div class="flex items-center gap-3">
                                <span class="text-sm text-gray-400">Every</span>
                                <input
                                    bind:value={intervalVal}
                                    type="number"
                                    min="1"
                                    class="w-20 bg-gray-900 border border-gray-700 rounded px-3 py-1.5 text-white text-center focus:ring-1 focus:ring-blue-500 outline-none"
                                />
                                <select
                                    bind:value={intervalUnit}
                                    class="bg-gray-900 border border-gray-700 rounded px-3 py-1.5 text-white focus:ring-1 focus:ring-blue-500 outline-none"
                                >
                                    <option value="m">Minute(s)</option>
                                    <option value="h">Hour(s)</option>
                                </select>
                            </div>
                        {:else}
                            <div>
                                <input
                                    bind:value={customCron}
                                    type="text"
                                    placeholder="0 10 * * *"
                                    class="w-full bg-gray-900 border border-gray-700 rounded px-4 py-2 text-white font-mono text-sm focus:ring-1 focus:ring-blue-500 outline-none"
                                />
                                <p class="text-xs text-gray-500 mt-2">
                                    Standard cron expression (min hour dom month
                                    dow)
                                </p>
                            </div>
                        {/if}
                    </div>
                </div>

                <div class="flex items-center gap-2">
                    <input
                        bind:checked={oneShot}
                        type="checkbox"
                        id="oneshot"
                        class="w-4 h-4 rounded border-gray-700 bg-gray-900 text-blue-600 focus:ring-offset-gray-900 cursor-pointer"
                    />
                    <label
                        for="oneshot"
                        class="text-sm text-gray-300 cursor-pointer"
                        >Run once and delete</label
                    >
                </div>
            </div>

            <div class="p-6 border-t border-gray-800 flex justify-end gap-3">
                <button
                    onclick={() => (isModalOpen = false)}
                    class="px-4 py-2 text-sm font-medium text-gray-400 hover:text-white transition-colors cursor-pointer"
                    >Cancel</button
                >
                <button
                    onclick={createTask}
                    disabled={isCreating || !name}
                    class="bg-blue-600 hover:bg-blue-700 disabled:bg-blue-600/50 disabled:cursor-not-allowed text-white px-6 py-2 rounded-lg font-medium transition-all shadow-lg hover:shadow-blue-900/20 active:scale-95 flex items-center gap-2 cursor-pointer"
                >
                    {#if isCreating}
                        <div
                            class="w-4 h-4 border-2 border-white/50 border-t-white rounded-full animate-spin"
                        ></div>
                    {:else}
                        <Save size={16} />
                    {/if}
                    Create Task
                </button>
            </div>
        </div>
    </div>
{/if}
