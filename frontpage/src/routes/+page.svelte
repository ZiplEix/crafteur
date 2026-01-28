<script lang="ts">
    import {
        ArrowRight,
        ChevronRight,
        Play,
        Layers,
        Activity,
        FolderTree,
        Package,
        Archive,
        Code,
        Copy,
        Check,
        Globe,
        Users,
        Zap,
        Calendar,
        ChevronLeft,
    } from "lucide-svelte";

    let copied = $state(false);
    const installCommand = "curl -sL get.crafteur.fr | sudo bash";

    function copyToClipboard() {
        navigator.clipboard.writeText(installCommand);
        copied = true;
        setTimeout(() => {
            copied = false;
        }, 2000);
    }

    const features = [
        {
            title: "Multi-Version",
            description:
                "Full support for Vanilla, Fabric (Mods), and Paper (Plugins). Deploy any version in just a few clicks.",
            icon: Layers,
            color: "text-green-400",
        },
        {
            title: "Real-Time Monitoring",
            description:
                "Monitor CPU, RAM, and performance live. Detailed graphs and automatic alerts.",
            icon: Activity,
            color: "text-blue-400",
        },
        {
            title: "File Manager",
            description:
                "Upload, download, ZIP/UNZIP compression, and file editing directly from the interface.",
            icon: FolderTree,
            color: "text-purple-400",
        },
        {
            title: "Integrated Modrinth",
            description:
                "Install your favorite mods and plugins in 1 click from the Modrinth catalog.",
            icon: Package,
            color: "text-emerald-400",
        },
        {
            title: "Automatic Backup",
            description:
                "Scheduled and compressed backups. One-click restoration to protect your worlds.",
            icon: Archive,
            color: "text-cyan-400",
        },
        {
            title: "Open Source",
            description:
                "100% free and self-hostable. Open code, active community, no limits.",
            icon: Code,
            color: "text-orange-400",
        },
        {
            title: "World Management",
            description:
                "Create multiple worlds, import existing maps (ZIP), and switch universes on the fly.",
            icon: Globe,
            color: "text-indigo-400",
        },
        {
            title: "Player Management",
            description:
                "Easily manage Whitelist, Operators (OP), and ban unwanted players.",
            icon: Users,
            color: "text-red-400",
        },
        {
            title: "Maximum Performance",
            description:
                "Powered by an ultra-lightweight Go backend for minimal resource usage.",
            icon: Zap,
            color: "text-yellow-400",
        },
        {
            title: "Automation",
            description:
                "Schedule automatic restarts, recurrent commands, and announcements.",
            icon: Calendar,
            color: "text-purple-400",
        },
    ];

    // Carousel Logic
    let currentSlide = $state(0);
    const screens = [
        {
            title: "Dashboard",
            src: "/mockup/mockup_dashboard_panel.png",
        },
        {
            title: "Console",
            src: "/mockup/mockup_console.png",
        },
        {
            title: "Schedules",
            src: "/mockup/mockup_schedule.png",
        },
        {
            title: "Backups",
            src: "/mockup/mockup_backup.png",
        },
        {
            title: "Files",
            src: "/mockup/mockup_save.png",
        },
        {
            title: "Worlds",
            src: "/mockup/mockup_world.png",
        },
        {
            title: "Add-ons",
            src: "/mockup/mockup_addson.png",
        },
        {
            title: "Configuration",
            src: "/mockup/mockup_configuration.png",
        },
        {
            title: "Players",
            src: "/mockup/mockup_player.png",
        },
    ];

    function nextSlide() {
        currentSlide = (currentSlide + 1) % screens.length;
    }

    function prevSlide() {
        currentSlide = (currentSlide - 1 + screens.length) % screens.length;
    }
</script>

<div class="relative isolate overflow-hidden">
    <!-- Background effects -->
    <div
        class="absolute inset-0 -z-10 bg-[radial-gradient(45rem_50rem_at_top,var(--color-slate-900),var(--color-slate-950))] opacity-20"
    ></div>
    <div
        class="absolute inset-y-0 right-1/2 -z-10 mr-16 w-[200%] origin-bottom-left skew-x-[-30deg] bg-slate-950 shadow-xl shadow-slate-900/10 ring-1 ring-slate-900 sm:mr-28 lg:mr-0 xl:mr-16 xl:origin-center"
    ></div>

    <div class="mx-auto max-w-7xl px-6 py-20 lg:px-8 lg:py-32">
        <div class="mx-auto max-w-2xl text-center">
            <h1
                class="text-4xl font-bold tracking-tight text-white sm:text-6xl"
            >
                Manage your servers <span
                    class="bg-clip-text text-transparent bg-linear-to-r from-green-400 to-blue-500"
                    >simply</span
                >
            </h1>
            <p class="mt-6 text-lg leading-8 text-slate-300">
                The open-source tool to deploy, manage, and monitor your
                Minecraft, Fabric, and Paper instances in seconds.
            </p>
            <div class="mt-10 flex items-center justify-center gap-x-6">
                <a
                    href="#install"
                    class="rounded-lg bg-green-600 px-6 py-3 text-sm font-semibold text-white shadow-xl shadow-green-900/20 hover:bg-green-500 hover:shadow-green-900/40 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-green-600 transition-all"
                >
                    Start Installation
                </a>
                <!-- <a
                    href="#demo"
                    class="group text-sm font-semibold leading-6 text-white hover:text-green-400 transition-colors flex items-center gap-2"
                >
                    Voir la Démo <span
                        aria-hidden="true"
                        class="group-hover:translate-x-1 transition-transform"
                        >→</span
                    >
                </a> -->
            </div>
        </div>

        <!-- 3D Mockup Placeholder -->
        <div class="mt-16 flow-root sm:mt-24">
            <div
                class="-m-2 rounded-xl bg-slate-900/50 ring-1 ring-inset ring-slate-900/10 lg:-m-4 lg:rounded-2xl lg:p-4 perspective-1000"
            >
                <div
                    class="relative rounded-xl bg-slate-800 shadow-2xl ring-1 ring-slate-900/10 aspect-video flex items-center justify-center transform hover:scale-[1.02] transition-transform duration-500 overflow-hidden group"
                    style="transform: rotateX(5deg);"
                >
                    <!-- Carousel Slide -->
                    {#each screens as screen, i}
                        <div
                            class="absolute inset-0 transition-opacity duration-500 flex items-center justify-center bg-slate-900"
                            class:opacity-100={currentSlide === i}
                            class:opacity-0={currentSlide !== i}
                        >
                            <img
                                src={screen.src}
                                alt={screen.title}
                                class="w-full h-full object-contain"
                            />

                            <!-- Overlay Gradient -->
                            <div
                                class="absolute inset-0 bg-gradient-to-t from-slate-950/20 to-transparent pointer-events-none"
                            ></div>
                        </div>
                    {/each}

                    <!-- Controls -->
                    <button
                        onclick={prevSlide}
                        aria-label="Previous slide"
                        class="absolute left-4 p-2 rounded-full bg-black/50 text-white opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/70 hover:cursor-pointer"
                    >
                        <ChevronLeft class="w-6 h-6" />
                    </button>
                    <button
                        onclick={nextSlide}
                        aria-label="Next slide"
                        class="absolute right-4 p-2 rounded-full bg-black/50 text-white opacity-0 group-hover:opacity-100 transition-opacity hover:bg-black/70 hover:cursor-pointer"
                    >
                        <ChevronRight class="w-6 h-6" />
                    </button>

                    <!-- Indicators -->
                    <div class="absolute bottom-4 flex gap-2">
                        {#each screens as _, i}
                            <button
                                aria-label="Go to slide {i + 1}"
                                onclick={() => (currentSlide = i)}
                                class="w-2 h-2 rounded-full transition-all hover:cursor-pointer {currentSlide ===
                                i
                                    ? 'bg-white w-4'
                                    : 'bg-white/50 hover:bg-white/80'}"
                            ></button>
                        {/each}
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>

<!-- Features Section -->
<section id="features" class="py-20 lg:py-32 bg-slate-950">
    <div class="mx-auto max-w-7xl px-6 lg:px-8">
        <div class="mx-auto max-w-2xl text-center mb-16">
            <h2
                class="text-3xl font-bold tracking-tight text-white sm:text-4xl"
            >
                Everything you need
            </h2>
            <p class="mt-4 text-lg text-slate-400">
                A complete suite of tools to manage your Minecraft servers like
                a pro.
            </p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each features as feature}
                <div
                    class="group relative bg-slate-900 border border-slate-800 rounded-xl p-6 hover:border-slate-700 transition-all duration-300 hover:shadow-xl hover:shadow-slate-900/50"
                >
                    <div class="flex items-start gap-4">
                        <div class="shrink-0">
                            <feature.icon class="h-8 w-8 {feature.color}" />
                        </div>
                        <div class="flex-1">
                            <h3 class="text-lg font-semibold text-white mb-2">
                                {feature.title}
                            </h3>
                            <p class="text-sm text-slate-400 leading-relaxed">
                                {feature.description}
                            </p>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    </div>
</section>

<!-- Installation Section -->
<section id="install" class="py-20 lg:py-32 bg-slate-900/50">
    <div class="mx-auto max-w-4xl px-6 lg:px-8">
        <div class="text-center mb-12">
            <h2
                class="text-3xl font-bold tracking-tight text-white sm:text-4xl"
            >
                Ready in 30 seconds
            </h2>
            <p class="mt-4 text-lg text-slate-400">
                A single command to install everything on your Ubuntu/Debian
                VPS.
            </p>
        </div>

        <!-- Terminal -->
        <div
            class="relative rounded-xl bg-black shadow-2xl ring-1 ring-slate-800"
        >
            <!-- Terminal Header -->
            <div
                class="flex items-center gap-2 border-b border-slate-800 px-4 py-3"
            >
                <div class="h-3 w-3 rounded-full bg-red-500"></div>
                <div class="h-3 w-3 rounded-full bg-yellow-500"></div>
                <div class="h-3 w-3 rounded-full bg-green-500"></div>
            </div>

            <!-- Terminal Content -->
            <div class="relative p-6">
                <pre
                    class="text-sm text-green-400 font-mono overflow-x-auto"><code
                        >{installCommand}</code
                    ></pre>

                <!-- Copy Button -->
                <button
                    onclick={copyToClipboard}
                    class="absolute top-4 right-4 rounded-lg bg-slate-800 p-2 text-slate-400 hover:bg-slate-700 hover:text-white transition-all"
                    aria-label="Copy to clipboard"
                >
                    {#if copied}
                        <Check class="h-5 w-5 text-green-400" />
                    {:else}
                        <Copy class="h-5 w-5" />
                    {/if}
                </button>
            </div>
        </div>

        {#if copied}
            <p class="mt-4 text-center text-sm text-green-400 animate-pulse">
                ✓ Copied to clipboard!
            </p>
        {/if}
    </div>
</section>
