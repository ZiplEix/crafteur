<script lang="ts">
    import { Github } from "lucide-svelte";
</script>

<footer class="border-t border-slate-800 bg-slate-950">
    <div class="mx-auto max-w-7xl px-6 py-12 lg:px-8 lg:py-16">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
            <!-- Logo & Copyright -->
            <div class="md:col-span-1">
                <a href="/" class="text-xl font-bold text-white tracking-tight">
                    Crafteur
                </a>
                <p class="mt-4 text-sm text-slate-400">
                    L'outil open-source pour gérer vos serveurs Minecraft.
                </p>
                <p class="mt-4 text-xs text-slate-500">
                    © 2026 Crafteur. Tous droits réservés.
                </p>
            </div>

            <!-- Produit -->
            <div>
                <h3 class="text-sm font-semibold text-white mb-4">Produit</h3>
                <ul class="space-y-3">
                    <li>
                        <a
                            href="#features"
                            class="text-sm text-slate-400 hover:text-white transition-colors"
                        >
                            Fonctionnalités
                        </a>
                    </li>
                    <li>
                        <a
                            href="#changelog"
                            class="text-sm text-slate-400 hover:text-white transition-colors"
                        >
                            Changelog
                        </a>
                    </li>
                </ul>
            </div>

            <!-- Ressources -->
            <div>
                <h3 class="text-sm font-semibold text-white mb-4">
                    Ressources
                </h3>
                <ul class="space-y-3">
                    <li>
                        <a
                            href="/docs"
                            class="text-sm text-slate-400 hover:text-white transition-colors"
                        >
                            Documentation
                        </a>
                    </li>
                    <li>
                        <a
                            href="https://github.com/ZiplEix/crafteur"
                            target="_blank"
                            rel="noreferrer"
                            class="text-sm text-slate-400 hover:text-white transition-colors flex items-center gap-1"
                        >
                            <Github class="h-4 w-4" />
                            GitHub
                        </a>
                    </li>
                    <li>
                        <a
                            href="#discord"
                            class="text-sm text-slate-400 hover:text-white transition-colors"
                        >
                            Discord
                        </a>
                    </li>
                </ul>
            </div>

            <!-- Légal -->
            <div>
                <h3 class="text-sm font-semibold text-white mb-4">Légal</h3>
                <ul class="space-y-3">
                    <li>
                        <a
                            href="https://github.com/ZiplEix/crafteur/blob/main/LICENSE"
                            target="_blank"
                            rel="noreferrer"
                            class="text-sm text-slate-400 hover:text-white transition-colors"
                        >
                            Licence MIT
                        </a>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</footer>
