---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
</script>

# Architecture des Fichiers

Comprendre l'organisation des fichiers de Crafteur est essentiel pour la maintenance et les sauvegardes.

## Arborescence

Crafteur installe toutes ses données dans `/opt/crafteur`. Voici la structure détaillée :

```text
/opt/crafteur/
├── crafteur                 # Le binaire exécutable de l'application
└── data/                    # Dossier de données persistantes
    ├── crafteur.db          # Base de données SQLite (utilisateurs, configs)
    ├── backups/             # Archives ZIP des sauvegardes serveurs
    └── servers/             # Dossiers des instances Minecraft
        └── <uuid>/          # Un dossier par serveur (nommé par UUID)
            ├── server.jar   # Jar du serveur (paper.jar, fabric.jar...)
            ├── server.properties
            ├── eula.txt
            ├── world/       # Données du monde
            ├── mods/        # Dossier mods (si Fabric)
            └── plugins/     # Dossier plugins (si Paper)
```

## Base de Données

Crafteur utilise **SQLite** pour stocker les métadonnées (comptes utilisateurs, configurations des serveurs, historique des tâches). Le fichier `crafteur.db` est le cœur de votre installation.

<Alert type="warning" title="Attention">
  Ne modifiez jamais le fichier <code>crafteur.db</code> ou le dossier <code>servers/</code> manuellement pendant que Crafteur est en cours d'exécution. Vous risquez de corrompre les données ou de provoquer des incohérences.
</Alert>

## Sauvegardes Manuelles

Bien que Crafteur dispose d'un systéme de backup automatique pour les serveurs Minecraft, il est recommandé de sauvegarder l'installation Crafteur elle-même régulièrement.

Pour faire un backup complet (Base de données + Serveurs) :

1. **Arrêtez le service** pour déverrouiller la base de données :
   ```bash
   sudo systemctl stop crafteur
   ```
2. **Copiez les données** :
   ```bash
   sudo cp -r /opt/crafteur/data /path/to/backup/location
   ```
3. **Redémarrez le service** :
   ```bash
   sudo systemctl start crafteur
   ```

<Alert type="error" title="Important">
  Assurez-vous toujours que le service est arrêté avant de copier `crafteur.db` pour éviter les fichiers corrompus.
</Alert>
