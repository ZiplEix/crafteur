export const configGroups = [
    {
        title: "Gameplay",
        fields: [
            { key: 'gamemode', label: 'Mode de jeu', type: 'select', options: ['survival', 'creative', 'adventure', 'spectator'] },
            { key: 'difficulty', label: 'Difficulté', type: 'select', options: ['peaceful', 'easy', 'normal', 'hard'] },
            { key: 'hardcore', label: 'Hardcore', type: 'boolean' },
            { key: 'pvp', label: 'PvP', type: 'boolean' },
            { key: 'allow-flight', label: 'Vol autorisé', type: 'boolean' },
        ]
    },
    {
        title: "Général",
        fields: [
            { key: 'motd', label: 'Message du jour (MOTD)', type: 'text' },
            { key: 'max-players', label: 'Joueurs Max', type: 'number' },
            { key: 'white-list', label: 'Whitelist', type: 'boolean' },
            { key: 'online-mode', label: 'Mode Online (Premium)', type: 'boolean' },
            { key: 'enable-command-block', label: 'Command Blocks', type: 'boolean' },
        ]
    },
    {
        title: "Monde & Performance",
        fields: [
            { key: 'level-seed', label: 'Graine (Seed)', type: 'text' },
            { key: 'view-distance', label: 'Distance de vue', type: 'number', min: 2, max: 32 },
            { key: 'simulation-distance', label: 'Distance de simulation', type: 'number', min: 2, max: 32 },
            { key: 'spawn-protection', label: 'Protection du Spawn', type: 'number' },
        ]
    }
];
