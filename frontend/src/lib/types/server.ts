export interface ServerStats {
    cpu: number;
    ram: number;
    ram_max: number;
}

export type WSMessageType = 'log' | 'status' | 'stats';

export interface WSMessage {
    type: WSMessageType;
    data: any;
}
