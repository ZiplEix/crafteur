<script lang="ts">
    import { Zap, Box, Settings, Github, Hammer, LogIn } from "lucide-svelte";
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
            <div class="grid grid-cols-1 gap-8 md:grid-cols-3">
                <!-- Feature 1 -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-8 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-green-500/10 p-3 text-green-500"
                    >
                        <Zap class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-xl font-bold text-white">
                        Ultra-fast Go Backend
                    </h3>
                    <p class="text-slate-400">
                        Designed for performance with a lightweight and reactive
                        backend written in Go.
                    </p>
                </div>

                <!-- Feature 2 -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-8 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-green-500/10 p-3 text-green-500"
                    >
                        <Box class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-xl font-bold text-white">
                        Vanilla & Fabric
                    </h3>
                    <p class="text-slate-400">
                        Native support for classic and modded Minecraft servers
                        with Fabric.
                    </p>
                </div>

                <!-- Feature 3 -->
                <div
                    class="group relative overflow-hidden rounded-xl border border-slate-800 bg-slate-900 p-8 transition hover:border-green-500/50 hover:shadow-lg hover:shadow-green-900/10"
                >
                    <div
                        class="mb-4 inline-flex items-center justify-center rounded-lg bg-green-500/10 p-3 text-green-500"
                    >
                        <Settings class="h-6 w-6" />
                    </div>
                    <h3 class="mb-2 text-xl font-bold text-white">
                        Visual Config
                    </h3>
                    <p class="text-slate-400">
                        Edit your server properties via an intuitive interface.
                    </p>
                </div>
            </div>
        </div>
    </section>
</div>
