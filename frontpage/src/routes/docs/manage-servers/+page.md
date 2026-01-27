---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
  let title = "Gestion des Serveurs";
</script>

# Gestion des Serveurs

Crafteur centralise toutes les opérations de maintenance de vos instances Minecraft. Cette page détaille les fonctionnalités du Dashboard.

## Cycle de Vie et Console

Le header du serveur vous permet de contrôler son état :
* **Démarrer :** Lance le processus Java avec les arguments de RAM définis.
* **Redémarrer :** Envoie une commande `stop`, attend l'arrêt du processus, puis le relance.
* **Arrêter :** Envoie la commande `stop` pour sauvegarder le monde proprement avant de couper le processus.

<Alert type="info" title="Kill vs Stop">
  Si le serveur ne répond plus (freeze), Crafteur détectera l'arrêt du processus après un délai. Il n'y a pas encore de bouton "Force Kill" dans l'interface V1 pour éviter la corruption de données.
</Alert>

## Configuration (`server.properties`)

L'onglet **Configuration** offre une interface graphique pour modifier le fichier `server.properties`.
* **Gameplay :** Mode de jeu (Survival/Creative), Difficulté, PvP.
* **Accès :** Whitelist, Nombre max de joueurs.
* **Réseau :** Port du serveur (par défaut 25565).

Chaque modification nécessite un **redémarrage** pour être prise en compte.

## Gestion des Versions

Vous pouvez changer la version de Minecraft (ex: passer de 1.20.1 à 1.20.4) ou le type de serveur (Vanilla ↔ Fabric).

<Alert type="error" title="Attention au Downgrade">
  Passer d'une version récente à une version plus ancienne (ex: 1.20 → 1.16) est impossible sans corrompre la carte. Minecraft ne sait pas "oublier" les nouveaux blocs. Faites toujours une sauvegarde avant de changer de version.
</Alert>

## Gestionnaire de Mondes (World Switcher)

L'onglet **Worlds** vous permet de gérer plusieurs dossiers de sauvegarde (maps) pour un même serveur.
1. **Création :** Générez un nouveau dossier (ex: `lobby`, `build_creative`).
2. **Activation :** Cliquez sur "Activer". Crafteur modifie la ligne `level-name` dans les propriétés.
3. **Application :** Redémarrez le serveur. Il chargera le nouveau dossier.

Ceci est idéal pour alterner entre une map "Survie" et une map "Mini-jeu" sans perdre de données.
