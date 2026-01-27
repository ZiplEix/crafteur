<script lang="ts">
    import {
        Zap,
        Github,
        Puzzle,
        Globe,
        Calendar,
        Users,
        Download,
        HardDrive,
        Layers,
        Activity,
        Code,
    } from "lucide-svelte";
</script>

<div class="flex flex-col">
    <!-- Hero Section -->
    <section
        class="mx-auto flex max-w-7xl flex-col items-center justify-center px-4 py-24 text-center sm:px-6 lg:px-8"
    >
        <h1
            class="text-4xl font-extrabold tracking-tight text-white sm:text-5xl md:text-6xl"
        >
            Your Minecraft server,<br />
            <span class="text-green-500">hassle-free.</span>
        </h1>
        <p class="mt-6 max-w-2xl text-lg text-slate-400">
            A modern and high-performance interface powered by a Go backend.
            Manage your Vanilla and Fabric servers with ease.
        </p>
        <div class="mt-10 flex gap-4">
            <a
                href="/dashboard"
                class="rounded-md bg-green-600 px-8 py-3 text-base font-medium text-white shadow-lg shadow-green-900/20 transition hover:bg-green-700 hover:shadow-green-900/40"
            >
                Access Dashboard
            </a>
            <a
                href="https://github.com/ZiplEix/crafteur"
                target="_blank"
                rel="noreferrer"
                class="flex items-center gap-2 rounded-md border border-slate-700 bg-slate-800/50 px-8 py-3 text-base font-medium text-slate-300 transition hover:bg-slate-800 hover:text-white"
            >
                <Github class="h-5 w-5" />
                <span>View Code</span>
            </a>
        </div>
    </section>

    <!-- Features Section -->
    <section class="w-full border-y border-slate-800 bg-slate-800/50 py-24">
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
            <div class="mb-16 text-center">
                <h2 class="text-3xl font-extrabold text-white sm:text-4xl">
                    Everything you need to run a server.
                </h2>
                <p class="mt-4 text-lg text-slate-400">
                    From simple vanilla survival to complex modded networks.
                </p>
            </div>

            <div class="grid grid-cols-1 gap-8 md:grid-cols-2 lg:grid-cols-4">
                <!-- Feature 1: Tech -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-green-500/10 p-3 text-green-500"
                    >
                        <Zap class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Blazing Fast
                    </h3>
                    <p class="text-sm text-slate-400">
                        Powered by a lightweight Go backend for minimal resource
                        usage and instant response times.
                    </p>
                </div>

                <!-- Feature 2: Multi-Loader -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-blue-500/10 p-3 text-blue-500"
                    >
                        <Layers class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Any Server Type
                    </h3>
                    <p class="text-sm text-slate-400">
                        Native support for <span class="text-white"
                            >Vanilla</span
                        >,
                        <span class="text-white">Fabric</span> (Mods), and
                        <span class="text-white">Paper</span> (Plugins).
                    </p>
                </div>

                <!-- Feature 3: Modrinth -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-emerald-500/10 p-3 text-emerald-500"
                    >
                        <Puzzle class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Add-on Browser
                    </h3>
                    <p class="text-sm text-slate-400">
                        Search and install Mods, Plugins, and Datapacks directly
                        from <span class="text-white">Modrinth</span>.
                    </p>
                </div>

                <!-- Feature 4: Files -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-yellow-500/10 p-3 text-yellow-500"
                    >
                        <HardDrive class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        File Manager
                    </h3>
                    <p class="text-sm text-slate-400">
                        Full web-based file browser. Upload, unzip, edit
                        configs, and manage files without FTP.
                    </p>
                </div>

                <!-- Feature 5: Automation -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-purple-500/10 p-3 text-purple-500"
                    >
                        <Calendar class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Automation
                    </h3>
                    <p class="text-sm text-slate-400">
                        Schedule automatic restarts, recurring commands, and
                        announcements easily.
                    </p>
                </div>

                <!-- Feature 6: Worlds -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-indigo-500/10 p-3 text-indigo-500"
                    >
                        <Globe class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        World Management
                    </h3>
                    <p class="text-sm text-slate-400">
                        Create multiple worlds, import existing ones via .zip,
                        and switch between them instantly.
                    </p>
                </div>

                <!-- Feature 7: Players -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-red-500/10 p-3 text-red-500"
                    >
                        <Users class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Player Control
                    </h3>
                    <p class="text-sm text-slate-400">
                        Manage Whitelist, Ops, and bans. See who is online in
                        real-time.
                    </p>
                </div>

                <!-- Feature 8: Backups -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-cyan-500/10 p-3 text-cyan-500"
                    >
                        <Download class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        One-Click Backups
                    </h3>
                    <p class="text-sm text-slate-400">
                        Create instant snapshots of your server and restore them
                        anytime. Safety first.
                    </p>
                </div>

                <!-- Feature 9: Monitoring -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-blue-500/10 p-3 text-blue-500"
                    >
                        <Activity class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Real-time Monitoring
                    </h3>
                    <p class="text-sm text-slate-400">
                        Watch CPU, RAM usage, and performance metrics live via
                        detailed charts.
                    </p>
                </div>

                <!-- Feature 10: Open Source -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-6 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-orange-500/10 p-3 text-orange-500"
                    >
                        <Code class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-lg font-bold text-white">
                        Open Source
                    </h3>
                    <p class="text-sm text-slate-400">
                        100% free and self-hostable. No limits, active
                        community, transparent code.
                    </p>
                </div>
            </div>
        </div>
    </section>
</div>
