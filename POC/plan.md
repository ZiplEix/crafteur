# Roadmap : Gestionnaire de Serveur Minecraft (Go + HTMX)

## 📅 Phase 1 : Le "Core" (Moteur Backend)
**Objectif :** Le programme Go sait télécharger, lancer et contrôler un serveur Minecraft (CLI uniquement).
- [ ] Initialiser le module Go (`go mod init`).
- [ ] Créer la structure de dossiers (`/cmd`, `/internal`, `/data`).
- [ ] Créer la fonction de téléchargement du `server.jar` Vanilla.
- [ ] Implémenter le wrapper de processus (`exec.Command`).
- [ ] Capturer les logs (STDOUT) vers le terminal Go.
- [ ] Envoyer des commandes (STDIN) depuis le terminal Go vers Minecraft.

## 📅 Phase 2 : Interface Web & Console Live
**Objectif :** Contrôle via navigateur avec logs en temps réel.
- [ ] Monter le serveur HTTP (net/http ou framework léger).
- [ ] Créer l'UI de base (HTML + HTMX + TailwindCSS).
- [ ] Connecter les boutons Start/Stop/Restart au backend.
- [ ] Implémenter les WebSockets pour streamer la console Minecraft vers le web.

## 📅 Phase 3 : Persistance & Configuration
**Objectif :** Gérer plusieurs serveurs et leurs propriétés.
- [ ] Intégrer SQLite (stockage des infos serveurs : id, port, type, nom).
- [ ] Créer le parser pour `server.properties` (Lecture/Écriture).
- [ ] Générer le formulaire de config dynamiquement via le parser.
- [ ] Gérer l'allocation dynamique des ports (éviter les conflits).

## 📅 Phase 4 : Mods, Fichiers & Fabric
**Objectif :** Support avancé (Modding) et gestion de fichiers.
- [ ] Créer un explorateur de fichiers web simple (au moins pour `/mods` et `/plugins`).
- [ ] Ajouter l'upload de fichiers (drag & drop).
- [ ] Implémenter la logique d'installation Fabric (Loader + Jar).

## 📅 Phase 5 : Packaging & Déploiement
**Objectif :** Installation "One-Click" sur Proxmox.
- [ ] Créer un Dockerfile Multi-stage (Build Go + Runtime Java).
- [ ] Optimiser l'image Docker (Base Alpine ou Debian Slim).
- [ ] Créer un `docker-compose.yml` de production.
- [ ] Documenter la commande d'installation unique.
