export interface DocSection {
	id: string;
	title: string;
	content: string;
}

export const docSections: DocSection[] = [
	{
		id: 'intro',
		title: 'Bienvenue sur Crafteur',
		content: `
			<p class="mb-4 text-slate-300">Crafteur est votre panneau de gestion de serveurs Minecraft simplifié. Il vous permet de créer, gérer et surveiller vos instances Minecraft en quelques clics, que vous soyez un joueur Vanilla, un amateur de Mods (Fabric) ou un administrateur de communauté (Paper/Plugins).</p>
		`
	},
	{
		id: 'create',
		title: 'Créer votre premier serveur',
		content: `
			<p class="mb-4 text-slate-300">Cliquez sur le bouton <strong>"Créer un serveur"</strong> sur le Dashboard.</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li><strong>Nom :</strong> Donnez un nom unique à votre serveur.</li>
				<li><strong>Type :</strong>
					<ul class="list-disc list-inside ml-6 mt-1 text-slate-400">
						<li><em>Vanilla :</em> Le jeu de base, sans modifications.</li>
						<li><em>Fabric :</em> Pour installer des <strong>Mods</strong>.</li>
						<li><em>Paper :</em> Optimisé pour les performances et les <strong>Plugins</strong>.</li>
					</ul>
				</li>
				<li><strong>Version :</strong> Choisissez la version de Minecraft (ex: 1.20.4).</li>
				<li><strong>Import :</strong> Vous pouvez envoyer un fichier <code>.zip</code> contenant un monde ou un serveur existant pour le restaurer immédiatement.</li>
			</ul>
		`
	},
	{
		id: 'manage',
		title: 'Pilotage & Console',
		content: `
			<p class="mb-4 text-slate-300">La page de détail vous offre un contrôle total :</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li><strong>Actions :</strong> Démarrer, Arrêter ou Redémarrer le serveur via les boutons en haut à droite.</li>
				<li><strong>Console :</strong> L'onglet Console vous montre ce qui se passe en temps réel. Vous pouvez envoyer des commandes (ex: <code>/gamemode creative</code>) directement via la barre de saisie en bas.</li>
				<li><strong>Monitoring :</strong> Surveillez la consommation CPU et RAM en haut de page pour prévenir les lags.</li>
			</ul>
		`
	},
	{
		id: 'files-addons',
		title: 'Fichiers, Mods et Plugins',
		content: `
			<p class="mb-2 text-slate-300"><strong>Onglet File :</strong> Un explorateur complet. Vous pouvez créer des dossiers, uploader des fichiers, supprimer des éléments et dézipper des archives directement sur le serveur.</p>
			<p class="mb-2 text-slate-300 font-semibold">Onglet Add-ons (Mods & Plugins) :</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Utilisez le <strong>Catalogue</strong> pour rechercher et installer des Mods (si Fabric) ou des Plugins (si Paper) depuis Modrinth en un clic.</li>
				<li>Les fichiers sont automatiquement placés dans le bon dossier.</li>
				<li>Vous pouvez aussi uploader vos propres fichiers <code>.jar</code> manuellement.</li>
			</ul>
		`
	},
	{
		id: 'players-worlds',
		title: 'Joueurs et Mondes',
		content: `
			<p class="mb-2 text-slate-300 font-semibold">Onglet Player : <span class="font-normal">Gérez votre communauté.</span></p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Voyez qui est en ligne (pastille verte).</li>
				<li>Nommez des Opérateurs (OP) ou bannissez des joueurs gênants.</li>
				<li>Sanctionnez avec un motif (Kick/Ban) via l'interface.</li>
			</ul>
			<p class="mb-2 text-slate-300 font-semibold">Onglet Worlds :</p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Créez plusieurs mondes (ex: "Survie", "Creatif", "Lobby").</li>
				<li>Activez un monde et redémarrez le serveur pour changer de map sans perdre l'ancienne.</li>
			</ul>
		`
	},
	{
		id: 'config',
		title: 'Configuration & Automatisation',
		content: `
			<p class="mb-2 text-slate-300"><strong>Onglet Configuration :</strong> Modifiez le fichier <code>server.properties</code> (Difficulté, PvP, Whitelist, etc.) via un formulaire simple. C'est aussi ici que vous pouvez changer la version du serveur ou le supprimer (Zone de Danger).</p>
			<p class="mb-2 text-slate-300 font-semibold">Onglet Schedule (Tâches) : <span class="font-normal">Automatisez votre serveur.</span></p>
			<ul class="list-disc list-inside space-y-2 text-slate-300 mb-4">
				<li>Programmez des redémarrages automatiques tous les jours.</li>
				<li>Envoyez des messages automatiques aux joueurs.</li>
				<li>Créez des backups réguliers.</li>
			</ul>
		`
	},
	{
		id: 'backups',
		title: 'Système de Backup',
		content: `
			<p class="mb-4 text-slate-300">Dans l'onglet <strong>Save</strong>, créez des instantanés (Snapshots) de votre serveur complet en un clic.</p>
			<p class="text-slate-300">Vous pouvez ensuite télécharger le fichier <code>.zip</code> pour le stocker sur votre ordinateur en sécurité. Il est recommandé de faire une sauvegarde avant toute mise à jour ou changement de version.</p>
		`
	}
];
