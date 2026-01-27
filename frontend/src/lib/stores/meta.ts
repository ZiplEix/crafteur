import { writable } from 'svelte/store';
import { api } from '$lib/api';
import type { MojangVersion } from '$lib/types/meta';

export const versions = writable<MojangVersion[]>([]);
export const loadingVersions = writable(false);

let loaded = false;

export async function loadVersions() {
    if (loaded) return;

    loadingVersions.set(true);
    try {
        const res = await api.get('/api/meta/versions');
        if (Array.isArray(res.data)) {
            versions.set(res.data);
            loaded = true;
        }
    } catch (e) {
        console.error("Failed to load versions", e);
    } finally {
        loadingVersions.set(false);
    }
}
