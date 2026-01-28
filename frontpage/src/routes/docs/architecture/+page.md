---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
</script>

# File Architecture

Understanding the organization of Crafteur files is essential for maintenance and backups.

## Directory Structure

Crafteur installs all its data in `/opt/crafteur`. Here is the detailed structure:

```text
/opt/crafteur/
├── crafteur                 # The application executable binary
└── data/                    # Persistent data directory
    ├── crafteur.db          # SQLite database (users, configs)
    ├── backups/             # Server backup ZIP archives
    └── servers/             # Minecraft instance directories
        └── <uuid>/          # One folder per server (named by UUID)
            ├── server.jar   # Server Jar (paper.jar, fabric.jar...)
            ├── server.properties
            ├── eula.txt
            ├── world/       # World data
            ├── mods/        # Mods folder (if Fabric)
            └── plugins/     # Plugins folder (if Paper)
```

## Database

Crafteur uses **SQLite** to store metadata (user accounts, server configurations, task history). The `crafteur.db` file is the heart of your installation.

<Alert type="warning" title="Warning">
  Never modify the <code>crafteur.db</code> file or the <code>servers/</code> folder manually while Crafteur is running. You risk corrupting data or causing inconsistencies.
</Alert>

## Manual Backups

Although Crafteur has an automatic backup system for Minecraft servers, it is recommended to back up the Crafteur installation itself regularly.

To make a full backup (Database + Servers):

1. **Stop the service** to unlock the database:
   ```bash
   sudo systemctl stop crafteur
   ```
2. **Copy the data**:
   ```bash
   sudo cp -r /opt/crafteur/data /path/to/backup/location
   ```
3. **Restart the service**:
   ```bash
   sudo systemctl start crafteur
   ```

<Alert type="error" title="Important">
  Always make sure the service is stopped before copying `crafteur.db` to avoid corrupted files.
</Alert>
