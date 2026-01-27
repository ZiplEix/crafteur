---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
  let title = "Mods, Plugins & Datapacks";
</script>

# Extensions (Add-ons)

Crafteur simplifie l'installation de contenu additionnel grâce à l'intégration de l'API **Modrinth**.

## Comprendre les types de contenus

Il est crucial de choisir le bon **Type de Serveur** à la création pour supporter les bonnes extensions :

| Type d'Extension | Type Serveur Requis | Dossier cible | Description |
| :--- | :--- | :--- | :--- |
| **Mods** | Fabric (ou Forge) | `/mods` | Modifie le jeu en profondeur. Nécessite souvent d'être installé sur le serveur **ET** chez le joueur. |
| **Plugins** | Paper (ou Spigot) | `/plugins` | Ajoute des fonctionnalités serveur (permissions, économie). Rien à installer pour le joueur. |
| **Datapacks** | Tous (Vanilla inclus) | `world/datapacks` | Scripts et modifications légères intégrés à la sauvegarde du monde. |

## Le Navigateur Modrinth

Dans l'onglet **Add-ons**, le sous-onglet "Catalogue" vous permet de rechercher des mods ou plugins.
* **Filtrage Intelligent :** Crafteur filtre automatiquement les résultats pour ne montrer que ceux compatibles avec votre version (ex: 1.20.1) et votre loader (Fabric/Paper).
* **Installation :** Un clic sur "Installer" télécharge la dernière version stable compatible.

<Alert type="warning" title="Gestion des Dépendances">
  Actuellement, Crafteur n'installe pas automatiquement les dépendances. Si vous installez un mod comme <em>Sodium</em>, vérifiez s'il nécessite <em>Fabric API</em> et installez-le manuellement via la recherche.
</Alert>

## Upload Manuel

Si vous possédez un fichier `.jar` ou `.zip` spécifique (non présent sur Modrinth ou développement privé), vous pouvez le glisser-déposer dans la zone d'upload ou utiliser le bouton "Ajouter".
* **Support Bulk :** Vous pouvez sélectionner plusieurs fichiers à la fois pour un upload groupé.
