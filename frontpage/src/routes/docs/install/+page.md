---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
</script>

# Installation Technique

Ce guide détaille le processus d'installation de Crafteur sur un serveur Linux.

## Prérequis

Pour installer et exécuter Crafteur, vous avez besoin de :
- Une distribution Linux récente (Ubuntu 22.04+ ou Debian 12+ recommandé)
- `systemd` comme système d'init
- `curl` pour télécharger le script d'installation
- Accès root ou sudo

## Le Script d'Installation

La commande d'installation rapide exécute un script bash qui automatise le déploiement :

```bash
curl -sL https://get.crafteur.fr | sudo bash
```

### Actions effectuées par le script

1. **Création de l'utilisateur** : Un utilisateur système `crafteur` est créé pour exécuter le service de manière sécurisée.
2. **Structure des dossiers** : Le répertoire `/opt/crafteur` est créé et les permissions sont attribuées à l'utilisateur `crafteur`.
3. **Installation du binaire** : La dernière version stable est téléchargée dans `/opt/crafteur/crafteur`.
4. **Service Systemd** : Un service est configuré pour lancer l'application au démarrage.

## Configuration Systemd

Voici le contenu typique du fichier `/etc/systemd/system/crafteur.service` généré :

```ini
[Unit]
Description=Crafteur Minecraft Server Manager
After=network.target

[Service]
Type=simple
User=crafteur
WorkingDirectory=/opt/crafteur
ExecStart=/opt/crafteur/crafteur
Restart=always
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
```

## Dépannage

### Port 8080 déjà utilisé

Par défaut, Crafteur écoute sur le port 8080. Si ce port est déjà occupé par un autre service, vous pouvez le changer en modifiant la variable d'environnement dans le service systemd.

Editez le fichier de service :
```bash
sudo nano /etc/systemd/system/crafteur.service
```

Modifiez la ligne `Environment` :
```ini
Environment=PORT=3000
```

Puis rechargez et redémarrez :
```bash
sudo systemctl daemon-reload
sudo systemctl restart crafteur
```

<Alert type="info" title="Conseil">
  Si vous utilisez un pare-feu (ufw), n'oubliez pas d'autoriser le port choisi : <code>sudo ufw allow 8080</code>.
</Alert>

## Désinstallation

Si vous souhaitez retirer Crafteur de votre serveur, nous fournissons un script automatisé qui nettoie le service et les fichiers binaires.

```bash
curl -sL https://raw.githubusercontent.com/ZiplEix/crafteur/main/uninstall.sh | sudo bash
```

Lors de l'exécution, le script vous demandera si vous souhaitez :

- **Conserver vos données (Mondes, Backups, Configuration)** : Seul le logiciel est supprimé.
- **Tout supprimer** : Le dossier `/opt/crafteur` est intégralement effacé.

<Alert type="error" title="Attention">
  Si vous choisissez de tout supprimer, cette action est irréversible. Pensez à télécharger vos backups via l'interface web (Onglet Save) avant de lancer cette commande.
</Alert>
