export interface DocSection {
	id: string;
	title: string;
	content: string;
}

export const docSections: DocSection[] = [
	{
		id: 'intro',
		title: 'Welcome to Crafteur',
		content: `
			<p class="mb-4 text-slate-300">Crafteur is your simplified Minecraft server management panel. It allows you to create, manage, and monitor your Minecraft instances in just a few clicks, whether you are a Vanilla player, a Mod enthusiast (Fabric), or a community administrator (Paper/Plugins).</p>
		`
	},
	{
		id: 'create',
		title: 'Create Your First Server',
		content: `
			<p class="mb-4 text-slate-300">Click the <strong>"Create Server"</strong> button on the Dashboard.</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li><strong>Name:</strong> Give your server a unique name.</li>
				<li><strong>Type:</strong>
					<ul class="list-disc list-inside ml-6 mt-1 text-slate-400">
						<li><em>Vanilla:</em> The base game, unmodified.</li>
						<li><em>Fabric:</em> For installing <strong>Mods</strong>.</li>
						<li><em>Paper:</em> Optimized for performance and <strong>Plugins</strong>.</li>
					</ul>
				</li>
				<li><strong>Version:</strong> Choose the Minecraft version (e.g., 1.20.4).</li>
				<li><strong>Import:</strong> You can upload a <code>.zip</code> file containing an existing world or server to restore it immediately.</li>
			</ul>
		`
	},
	{
		id: 'manage',
		title: 'Control & Console',
		content: `
			<p class="mb-4 text-slate-300">The details page offers you full control:</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li><strong>Actions:</strong> Start, Stop, or Restart the server via the buttons at the top right.</li>
				<li><strong>Console:</strong> The Console tab shows what is happening in real-time. You can send commands (e.g., <code>/gamemode creative</code>) directly via the input bar at the bottom.</li>
				<li><strong>Monitoring:</strong> Monitor CPU and RAM usage at the top of the page to prevent lags.</li>
			</ul>
		`
	},
	{
		id: 'files-addons',
		title: 'Files, Mods & Plugins',
		content: `
			<p class="mb-2 text-slate-300"><strong>Files Tab:</strong> A complete explorer. You can create folders, upload files, delete items, and unzip archives directly on the server.</p>
			<p class="mb-2 text-slate-300 font-semibold">Add-ons Tab (Mods & Plugins):</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Use the <strong>Catalog</strong> to search and install Mods (if Fabric) or Plugins (if Paper) from Modrinth in one click.</li>
				<li>Files are automatically placed in the correct folder.</li>
				<li>You can also upload your own <code>.jar</code> files manually.</li>
			</ul>
		`
	},
	{
		id: 'players-worlds',
		title: 'Players & Worlds',
		content: `
			<p class="mb-2 text-slate-300 font-semibold">Players Tab: <span class="font-normal">Manage your community.</span></p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>See who is online (green dot).</li>
				<li>Appoint Operators (OP) or ban troublesome players.</li>
				<li>Sanction with a reason (Kick/Ban) via the interface.</li>
			</ul>
			<p class="mb-2 text-slate-300 font-semibold">Worlds Tab:</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Create multiple worlds (e.g., "Survival", "Creative", "Lobby").</li>
				<li>Activate a world and restart the server to change the map without losing the old one.</li>
			</ul>
		`
	},
	{
		id: 'config',
		title: 'Configuration & Automation',
		content: `
			<p class="mb-2 text-slate-300"><strong>Configuration Tab:</strong> Modify the <code>server.properties</code> file (Difficulty, PvP, Whitelist, etc.) via a simple form. This is also where you can change the server version or delete it (Danger Zone).</p>
			<p class="mb-2 text-slate-300 font-semibold">Schedules Tab: <span class="font-normal">Automate your server.</span></p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Schedule automatic restarts every day.</li>
				<li>Send automatic messages to players.</li>
				<li>Create regular backups.</li>
			</ul>
		`
	},
	{
		id: 'backups',
		title: 'Backup System',
		content: `
			<p class="mb-4 text-slate-300">In the <strong>Backups</strong> tab, create snapshots of your entire server in one click.</p>
			<p class="text-slate-300">You can then download the <code>.zip</code> file to store it safely on your computer. It is recommended to make a backup before any update or version change.</p>
		`
	}
];
