---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
</script>

# API Reference

Crafteur exposes a RESTful API that allows you to automate the management of your servers or integrate Crafteur with other tools.

## Authentication

All requests to the API must be authenticated via a Token.

1. Log in to the web interface.
2. Open developer tools (F12) or inspect requests to retrieve your `Authorization` token.
3. Use this token in the header of your requests:

```http
Authorization: Bearer <your_jwt_token>
```

<Alert type="info" title="Note">
  A feature to generate persistent API Keys is planned for future versions.
</Alert>

## Usage Examples

Here are some common examples of interacting with the API.

### Start a server

Allows launching a specific Minecraft instance.

- **Method**: `POST`
- **Endpoint**: `/api/servers/{id}/start`

**cURL Example:**

```bash
curl -X POST https://panel.yourdomain.com/api/servers/123e4567-e89b-12d3-a456-426614174000/start \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json"
```

### Stop a server

Properly stops an instance.

- **Method**: `POST`
- **Endpoint**: `/api/servers/{id}/stop`

### Get statuses

Get the list of servers and their status (online/offline).

- **Method**: `GET`
- **Endpoint**: `/api/servers`

Response (JSON):
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
