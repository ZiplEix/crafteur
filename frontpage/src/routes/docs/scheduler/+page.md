---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
  let title = "Automatisation & Tâches";
</script>

# Tâches Planifiées (Scheduler)

Le système de planification vous permet d'automatiser la maintenance de votre serveur. Il est basé sur le format standard **Cron**, mais propose une interface visuelle simplifiée.

## Types de Déclencheurs

### 1. Récurrent (Jours/Heures)
Idéal pour des actions à heure fixe.
* *Exemple :* "Tous les jours à 04:00 du matin".
* *Interface :* Sélectionnez l'heure et cochez les jours de la semaine souhaités.

### 2. Intervalle
Idéal pour des actions répétitives fréquentes.
* *Exemple :* "Toutes les 30 minutes".
* *Usage :* Messages automatiques, sauvegardes fréquentes.

### 3. Custom (Cron Avancé)
Pour les utilisateurs experts nécessitant une précision spécifique.
* *Format :* `Minute Heure Jour Mois JourSemaine`
* *Exemple :* `0 12 1 * *` (À midi, le 1er jour de chaque mois).

## Actions Disponibles

Une tâche peut déclencher une des actions suivantes :
* **Start / Stop / Restart :** Contrôle l'état du serveur.
* **Command :** Exécute une ou plusieurs commandes Minecraft dans la console.

<Alert type="success" title="Astuce : Chaînage de Commandes">
  <p>L'action "Command" supporte le multi-lignes. Vous pouvez créer une séquence complexe :</p>
  <pre class="text-xs bg-black p-2 rounded text-slate-300">
say Attention, redémarrage dans 1 minute !
save-all
stop
  </pre>
  <p>Couplez ceci avec une tâche "Start" programmée 2 minutes plus tard pour un cycle de maintenance complet.</p>
</Alert>

## Exemples Courants

| Objectif | Type | Fréquence | Action |
| :--- | :--- | :--- | :--- |
| **Reboot Quotidien** | Récurrent | 06:00 (Tous les jours) | Restart |
| **Message de Bienvenue** | Intervalle | 2 Heures | Command: `say Rejoignez notre Discord !` |
| **Sauvegarde Monde** | Intervalle | 1 Heure | Command: `save-all` |
