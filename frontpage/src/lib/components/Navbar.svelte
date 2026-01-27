<script lang="ts">
    import { Github, Download } from "lucide-svelte";
</script>

<nav
    class="sticky top-0 z-50 w-full border-b border-white/5 bg-slate-900/80 backdrop-blur-md"
>
    <div
        class="mx-auto flex h-16 max-w-7xl items-center justify-between px-4 sm:px-6 lg:px-8"
    >
        <!-- Left: Logo -->
        <div class="flex items-center">
            <a href="/" class="text-xl font-bold text-white tracking-tight">
                Crafteur
            </a>
        </div>

        <!-- Center: Links -->
        <div class="hidden md:flex items-center space-x-8">
            <a
                href="/#features"
                class="text-sm font-medium text-slate-300 hover:text-white transition-colors"
            >
                Fonctionnalités
            </a>
            <a
                href="/#install"
                class="text-sm font-medium text-slate-300 hover:text-white transition-colors"
            >
                Installation
            </a>
            <a
                href="/docs"
                class="text-sm font-medium text-slate-300 hover:text-white transition-colors"
            >
                Documentation
            </a>
        </div>

        <!-- Right: Buttons -->
        <div class="flex items-center space-x-4">
            <a
                href="https://github.com/ZiplEix/crafteur"
                target="_blank"
                rel="noreferrer"
                class="text-slate-400 hover:text-white transition-colors"
                aria-label="GitHub"
            >
                <Github class="h-6 w-6" />
            </a>

            <a
                href="#download"
                class="inline-flex items-center gap-2 rounded-lg bg-green-600 px-4 py-2 text-sm font-semibold text-white shadow-lg shadow-green-900/20 hover:bg-green-500 hover:shadow-green-900/40 transition-all focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-slate-900"
            >
                <Download class="h-4 w-4" />
                <span>Télécharger</span>
            </a>
        </div>
    </div>
</nav>
