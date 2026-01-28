---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
</script>

# Technical Installation

This guide details the installation process of Crafteur on a Linux server.

## Prerequisites

To install and run Crafteur, you need:
- A recent Linux distribution (Ubuntu 22.04+ or Debian 12+ recommended)
- `systemd` as the init system
- `curl` to download the installation script
- Root or sudo access

## The Installation Script

The quick installation command runs a bash script that automates the deployment:

```bash
curl -sL get.crafteur.fr | sudo bash
```

### Actions performed by the script

1. **User Creation**: A system user `crafteur` is created to run the service securely.
2. **Folder Structure**: The `/opt/crafteur` directory is created and permissions are assigned to the `crafteur` user.
3. **Binary Installation**: The latest stable version is downloaded to `/opt/crafteur/crafteur`.
4. **Systemd Service**: A service is configured to launch the application at startup.

## Systemd Configuration

Here is the typical content of the generated `/etc/systemd/system/crafteur.service` file:

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

## Troubleshooting

### Port 8080 already in use

By default, Crafteur listens on port 8080. If this port is already occupied by another service, you can change it by modifying the environment variable in the systemd service.

Edit the service file:
```bash
sudo nano /etc/systemd/system/crafteur.service
```

Modify the `Environment` line:
```ini
Environment=PORT=3000
```

Then reload and restart:
```bash
sudo systemctl daemon-reload
sudo systemctl restart crafteur
```

<Alert type="info" title="Tip">
  If you use a firewall (ufw), don't forget to allow the chosen port: <code>sudo ufw allow 8080</code>.
</Alert>

## Uninstallation

If you wish to remove Crafteur from your server, we provide an automated script that cleans up the service and binary files.

```bash
curl -sL https://raw.githubusercontent.com/ZiplEix/crafteur/main/uninstall.sh | sudo bash
```

During execution, the script will ask if you wish to:

- **Keep your data (Worlds, Backups, Configuration)**: Only the software is removed.
- **Delete everything**: The `/opt/crafteur` directory is completely erased.

<Alert type="error" title="Warning">
  If you choose to delete everything, this action is irreversible. Remember to download your backups via the web interface (Save Tab) before running this command.
</Alert>
