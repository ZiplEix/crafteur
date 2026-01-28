---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
  let title = "Managing Servers";
</script>

# Managing Servers

Crafteur centralizes all maintenance operations for your Minecraft instances. This page details the features of the Dashboard.

## Lifecycle and Console

The server header allows you to control its state:
* **Start:** Launches the Java process with the defined RAM arguments.
* **Restart:** Sends a `stop` command, waits for the process to stop, then relaunches it.
* **Stop:** Sends the `stop` command to save the world properly before cutting the process.

<Alert type="info" title="Kill vs Stop">
  If the server is no longer responding (freeze), Crafteur will detect the process stoppage after a delay. There is no "Force Kill" button yet in the V1 interface to avoid data corruption.
</Alert>

## Configuration (`server.properties`)

The **Configuration** tab offers a graphical interface to modify the `server.properties` file.
* **Gameplay:** Game Mode (Survival/Creative), Difficulty, PvP.
* **Access:** Whitelist, Max players.
* **Network:** Server port (default 25565).

Each modification requires a **restart** to take effect.

## Version Management

You can change the Minecraft version (e.g., from 1.20.1 to 1.20.4) or the server type (Vanilla ↔ Fabric).

<Alert type="error" title="Downgrade Warning">
  Downgrading from a recent version to an older version (e.g., 1.20 → 1.16) is impossible without corrupting the map. Minecraft does not know how to "forget" new blocks. Always make a backup before changing versions.
</Alert>

## World Switcher

The **Worlds** tab allows you to manage multiple save folders (maps) for the same server.
1. **Creation:** Generate a new folder (e.g., `lobby`, `build_creative`).
2. **Activation:** Click "Activate". Crafteur modifies the `level-name` line in the properties.
3. **Application:** Restart the server. It will load the new folder.

This is ideal for switching between a "Survival" map and a "Minigame" map without losing data.
