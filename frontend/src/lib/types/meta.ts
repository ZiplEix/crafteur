export interface MojangVersion {
    id: string;
    type: string; // "release" | "snapshot" | "old_beta" | "old_alpha"
    url: string;
    time?: string;
    releaseTime?: string;
}
