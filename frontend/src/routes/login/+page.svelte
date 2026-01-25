<script lang="ts">
    import { goto } from "$app/navigation";
    import { api } from "$lib/api";

    let username = $state("");
    let password = $state("");
    let error = $state("");

    async function handleSubmit(event: Event) {
        event.preventDefault();
        error = "";

        try {
            await api.post("/api/login", { username, password });
            goto("/dashboard");
        } catch (e: any) {
            console.error(e);
            if (e.response && e.response.data && e.response.data.error) {
                error = e.response.data.error;
            } else {
                error = "Erreur de connexion au serveur";
            }
        }
    }
</script>

<div class="flex min-h-[calc(100vh-200px)] items-center justify-center p-4">
    <div
        class="w-full max-w-md rounded-xl bg-slate-800 p-8 shadow-2xl border border-slate-700"
    >
        <h1 class="mb-6 text-center text-3xl font-bold text-white">
            Connexion
        </h1>

        {#if error}
            <div
                class="mb-4 rounded bg-red-500/10 p-3 text-sm text-red-500 border border-red-500/20 text-center"
            >
                {error}
            </div>
        {/if}

        <form onsubmit={handleSubmit} class="space-y-6">
            <div>
                <label
                    for="username"
                    class="mb-2 block text-sm font-medium text-slate-300"
                    >Nom d'utilisateur</label
                >
                <input
                    type="text"
                    id="username"
                    bind:value={username}
                    class="w-full rounded-lg bg-slate-900 border border-slate-700 p-3 text-white placeholder-slate-500 focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500 transition-colors"
                    placeholder="Votre nom d'utilisateur"
                    required
                />
            </div>

            <div>
                <label
                    for="password"
                    class="mb-2 block text-sm font-medium text-slate-300"
                    >Mot de passe</label
                >
                <input
                    type="password"
                    id="password"
                    bind:value={password}
                    class="w-full rounded-lg bg-slate-900 border border-slate-700 p-3 text-white placeholder-slate-500 focus:border-green-500 focus:outline-none focus:ring-1 focus:ring-green-500 transition-colors"
                    placeholder="••••••••"
                    required
                />
            </div>

            <button
                type="submit"
                class="w-full rounded-lg bg-green-600 px-4 py-3 font-semibold text-white hover:bg-green-500 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 focus:ring-offset-slate-800 transition-colors"
            >
                Se connecter
            </button>
        </form>
    </div>
</div>
