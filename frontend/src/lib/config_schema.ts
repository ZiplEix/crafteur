export const configGroups = [
    {
        title: "Gameplay",
        fields: [
            { key: 'gamemode', label: 'Gamemode', type: 'select', options: ['survival', 'creative', 'adventure', 'spectator'] },
            { key: 'difficulty', label: 'Difficulty', type: 'select', options: ['peaceful', 'easy', 'normal', 'hard'] },
            { key: 'hardcore', label: 'Hardcore', type: 'boolean' },
            { key: 'pvp', label: 'PvP', type: 'boolean' },
            { key: 'allow-flight', label: 'Allow Flight', type: 'boolean' },
        ]
    },
    {
        title: "General",
        fields: [
            { key: 'motd', label: 'Message of the Day (MOTD)', type: 'text' },
            { key: 'max-players', label: 'Max Players', type: 'number' },
            { key: 'white-list', label: 'Whitelist', type: 'boolean' },
            { key: 'online-mode', label: 'Online Mode (Premium)', type: 'boolean' },
            { key: 'enable-command-block', label: 'Command Blocks', type: 'boolean' },
        ]
    },
    {
        title: "World & Performance",
        fields: [
            { key: 'level-seed', label: 'Seed', type: 'text' },
            { key: 'view-distance', label: 'View Distance', type: 'number', min: 2, max: 32 },
            { key: 'simulation-distance', label: 'Simulation Distance', type: 'number', min: 2, max: 32 },
            { key: 'spawn-protection', label: 'Spawn Protection', type: 'number' },
        ]
    }
];
