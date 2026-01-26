export interface PlayerCache {
    name: string;
    uuid: string;
    expiresOn: string;
    online: boolean;
}

export interface OpEntry {
    uuid: string;
    name: string;
    level: number;
    bypassesPlayerLimit: boolean;
}

export interface BanEntry {
    uuid: string;
    name: string;
    created: string;
    source: string;
    expires: string;
    reason: string;
}

export interface PlayerActionRequest {
    player: string;
    action: 'op' | 'deop' | 'ban' | 'pardon' | 'kick' | 'whitelist_add' | 'whitelist_remove';
    reason?: string;
}
