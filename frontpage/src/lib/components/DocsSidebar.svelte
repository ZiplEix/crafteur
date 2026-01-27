<script lang="ts">
    import { page } from "$app/stores";
    import { docMenu } from "$lib/data/docsParams";

    // Helper to check if a link is active
    // We check if the current path starts with the link href to handle potential sub-pages correctly
    // But for exact matches in strict navigation, strict equality is often better.
    // Given the request spec for "Active Link", we'll check equality for precise highlighting.
    const isActive = (href: string) => $page.url.pathname === href;
</script>

<nav class="p-6 space-y-8">
    {#each docMenu as category}
        <div>
            <h3
                class="text-xs font-bold text-slate-500 uppercase tracking-wider mb-3"
            >
                {category.title}
            </h3>
            <ul class="space-y-2">
                {#each category.items as item}
                    <li>
                        <a
                            href={item.href}
                            class="block py-1 text-sm transition-colors {isActive(
                                item.href,
                            )
                                ? 'text-blue-400 font-medium border-l-2 border-blue-400 pl-3 -ml-3'
                                : 'text-slate-400 hover:text-white'}"
                        >
                            {item.title}
                        </a>
                    </li>
                {/each}
            </ul>
        </div>
    {/each}
</nav>
