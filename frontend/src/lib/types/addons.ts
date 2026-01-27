export interface AddonFile {
    name: string;
    is_dir: boolean;
    size: number;
    mod_time: string;
}

export type AddonType = 'mods' | 'plugins' | 'datapacks';
