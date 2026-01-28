---
layout: docs
---

<script>
  import Alert from '$lib/components/Alert.svelte';
  let title = "Automation & Tasks";
</script>

# Scheduled Tasks (Scheduler)

The scheduling system allows you to automate your server maintenance. It is based on the standard **Cron** format but offers a simplified visual interface.

## Trigger Types

### 1. Recurring (Days/Hours)
Ideal for actions at a fixed time.
* *Example:* "Every day at 04:00 AM".
* *Interface:* Select the time and check the desired days of the week.

### 2. Interval
Ideal for frequent repetitive actions.
* *Example:* "Every 30 minutes".
* *Usage:* Automatic messages, frequent backups.

### 3. Custom (Advanced Cron)
For expert users requiring specific precision.
* *Format:* `Minute Hour Day Month DayOfWeek`
* *Example:* `0 12 1 * *` (At noon, on the 1st day of every month).

## Available Actions

A task can trigger one of the following actions:
* **Start / Stop / Restart:** Controls the server state.
* **Command:** Executes one or more Minecraft commands in the console.

<Alert type="success" title="Tip: Command Chaining">
  <p>The "Command" action supports multi-lines. You can create a complex sequence:</p>
  <pre class="text-xs bg-black p-2 rounded text-slate-300">
say Attention, restart in 1 minute!
save-all
stop
  </pre>
  <p>Couple this with a "Start" task scheduled 2 minutes later for a full maintenance cycle.</p>
</Alert>

## Common Examples

| Objective | Type | Frequency | Action |
| :--- | :--- | :--- | :--- |
| **Daily Reboot** | Recurring | 06:00 (Every day) | Restart |
| **Welcome Message** | Interval | 2 Hours | Command: `say Join our Discord!` |
| **World Backup** | Interval | 1 Hour | Command: `save-all` |
