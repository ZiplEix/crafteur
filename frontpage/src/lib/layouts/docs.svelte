<script lang="ts">
    import DocsSidebar from "$lib/components/DocsSidebar.svelte";
    import { Menu, X } from "lucide-svelte";

    let { children } = $props();
    let isMobileMenuOpen = $state(false);

    function toggleMobileMenu() {
        isMobileMenuOpen = !isMobileMenuOpen;
    }
</script>

<div class="flex min-h-screen bg-slate-950 text-slate-200">
    <!-- Mobile Menu Button -->
    <button
        onclick={toggleMobileMenu}
        class="fixed top-4 left-4 z-50 rounded-lg bg-slate-800 p-2 text-slate-200 lg:hidden"
        aria-label="Toggle Menu"
    >
        {#if isMobileMenuOpen}
            <X class="h-6 w-6" />
        {:else}
            <Menu class="h-6 w-6" />
        {/if}
    </button>

    <!-- Sidebar Overlay (Mobile) -->
    {#if isMobileMenuOpen}
        <div
            class="fixed inset-0 z-40 bg-slate-950/80 backdrop-blur-sm lg:hidden"
            onclick={toggleMobileMenu}
            role="button"
            tabindex="0"
            aria-label="Close menu"
            onkeydown={(e) => e.key === "Escape" && toggleMobileMenu()}
        ></div>
    {/if}

    <!-- Sidebar -->
    <aside
        class="fixed inset-y-0 left-0 z-40 w-64 transform border-r border-slate-800 bg-slate-900 transition-transform duration-300 ease-in-out lg:translate-x-0 {isMobileMenuOpen
            ? 'translate-x-0'
            : '-translate-x-full'}"
    >
        <div class="h-full overflow-y-auto pt-20">
            <DocsSidebar />
        </div>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 lg:ml-64">
        <article
            class="prose prose-invert prose-blue max-w-4xl mx-auto py-12 px-6 pt-20 lg:pt-12"
        >
            {@render children()}
        </article>
    </div>
</div>
