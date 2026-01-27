---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
</script>

# Référence API

Crafteur expose une API RESTful qui vous permet d'automatiser la gestion de vos serveurs ou d'intégrer Crafteur à d'autres outils.

## Authentification

Toutes les requêtes vers l'API doivent être authentifiées via un Jeton (Token). 

1. Connectez-vous à l'interface web.
2. Ouvrez les outils de développement (F12) ou inspectez les requêtes pour récupérer votre token `Authorization`.
3. Utilisez ce token dans le header de vos requêtes :

```http
Authorization: Bearer <votre_token_jwt>
```

<Alert type="info" title="Note">
  Une fonctionnalité pour générer des clés API persistantes (API Keys) est prévue pour les prochaines versions.
</Alert>

## Exemples d'Utilisation

Voici quelques exemples courants d'interaction avec l'API.

### Démarrer un serveur

Permet de lancer une instance Minecraft spécifique.

- **Méthode** : `POST`
- **Endpoint** : `/api/servers/{id}/start`

**Exemple cURL :**

```bash
curl -X POST https://panel.votredomaine.com/api/servers/123e4567-e89b-12d3-a456-426614174000/start \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json"
```

### Arrêter un serveur

Arrête proprement une instance.

- **Méthode** : `POST`
- **Endpoint** : `/api/servers/{id}/stop`

### Récupérer les statuts

Obtenir la liste des serveurs et leur statut (online/offline).

- **Méthode** : `GET`
- **Endpoint** : `/api/servers`

Reponse (JSON) :
```json
[
  {
    "id": "123e4567-...",
    "name": "Survival Server",
    "status": "running",
    "port": 25565
  },
  {
    "id": "987fcdeb-...",
    "name": "Creative Plot",
    "status": "stopped",
    "port": 25566
  }
]
```
