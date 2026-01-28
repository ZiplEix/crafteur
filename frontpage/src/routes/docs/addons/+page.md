---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
  let title = "Mods, Plugins & Datapacks";
</script>

# Extensions (Add-ons)

Crafteur simplifies the installation of additional content thanks to the **Modrinth** API integration.

## Understanding Content Types

It is crucial to choose the correct **Server Type** at creation to support the right extensions:

| Extension Type | Required Server Type | Target Folder | Description |
| :--- | :--- | :--- | :--- |
| **Mods** | Fabric (or Forge) | `/mods` | Modifies the game deeply. Often requires installation on the server **AND** the player's client. |
| **Plugins** | Paper (or Spigot) | `/plugins` | Adds server features (permissions, economy). Nothing to install for the player. |
| **Datapacks** | All (Vanilla included) | `world/datapacks` | Scripts and light modifications integrated into the world save. |

## The Modrinth Browser

In the **Add-ons** tab, the "Catalog" sub-tab allows you to search for mods or plugins.
* **Smart Filtering:** Crafteur automatically filters results to show only those compatible with your version (e.g., 1.20.1) and your loader (Fabric/Paper).
* **Installation:** One click on "Install" downloads the latest stable compatible version.

<Alert type="warning" title="Dependency Management">
  Currently, Crafteur does not automatically install dependencies. If you install a mod like <em>Sodium</em>, check if it requires <em>Fabric API</em> and install it manually via search.
</Alert>

## Manual Upload

If you have a specific `.jar` or `.zip` file (not present on Modrinth or private development), you can drag and drop it into the upload area or use the "Add" button.
* **Bulk Support:** You can select multiple files at once for bulk upload.
