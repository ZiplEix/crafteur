<script lang="ts">
    import { docSections } from "$lib/data/docs";
    import { page } from "$app/stores";
    import { onMount } from "svelte";

    let activeSection = $state("intro");

    function scrollToSection(id: string) {
        const element = document.getElementById(id);
        if (element) {
            element.scrollIntoView({ behavior: "smooth" });
            activeSection = id;
            // Update URL hash without jumping
            history.pushState(null, "", `#${id}`);
        }
    }

    onMount(() => {
        const hash = $page.url.hash.substring(1);
        if (hash && docSections.some((s) => s.id === hash)) {
            scrollToSection(hash);
        }
    });
</script>

<div class="flex h-[calc(100vh-4rem)]">
    <!-- Sidebar -->
    <aside
        class="w-64 border-r border-slate-700 bg-slate-900/50 overflow-y-auto hidden md:block sticky top-0 h-full"
    >
        <nav class="p-4 space-y-1">
            <h3
                class="px-3 text-xs font-semibold text-slate-500 uppercase tracking-wider mb-2"
            >
                Documentation
            </h3>
            {#each docSections as section}
                <button
                    onclick={() => scrollToSection(section.id)}
                    class="w-full text-left rounded-md px-3 py-2 text-sm font-medium transition-colors {activeSection ===
                    section.id
                        ? 'bg-slate-800 text-white'
                        : 'text-slate-400 hover:bg-slate-800 hover:text-white'}"
                >
                    {section.title}
                </button>
            {/each}
        </nav>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 overflow-y-auto bg-slate-900">
        <div class="max-w-4xl mx-auto px-8 py-12">
            <div class="space-y-12">
                {#each docSections as section}
                    <section id={section.id} class="scroll-mt-20">
                        <h2
                            class="text-3xl font-bold text-white mb-6 pb-2 border-b border-slate-800"
                        >
                            {section.title}
                        </h2>
                        <div class="text-slate-300 leading-relaxed space-y-4">
                            {@html section.content}
                        </div>
                    </section>
                {/each}
            </div>

            <!-- Quick footer for docs -->
            <div
                class="mt-20 pt-8 border-t border-slate-800 text-center text-slate-500 text-sm"
            >
                <p>Besoin de plus d'aide ? Contactez le support.</p>
            </div>
        </div>
    </div>
</div>
