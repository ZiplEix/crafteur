package database

import (
	"database/sql"
	"time"

	"github.com/ZiplEix/crafteur/core"
)

func CreateTask(task *core.ScheduledTask) error {
	stmt, err := DB.Prepare(`INSERT INTO scheduled_tasks (id, server_id, name, action, payload, cron_expression, one_shot, last_run) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(task.ID, task.ServerID, task.Name, task.Action, task.Payload, task.CronExpression, task.OneShot, task.LastRun)
	return err
}

func GetTasksByServer(serverID string) ([]core.ScheduledTask, error) {
	rows, err := DB.Query(`SELECT id, server_id, name, action, payload, cron_expression, one_shot, last_run FROM scheduled_tasks WHERE server_id = ?`, serverID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []core.ScheduledTask
	for rows.Next() {
		var t core.ScheduledTask
		var lastRun sql.NullTime
		if err := rows.Scan(&t.ID, &t.ServerID, &t.Name, &t.Action, &t.Payload, &t.CronExpression, &t.OneShot, &lastRun); err != nil {
			return nil, err
		}
		if lastRun.Valid {
			t.LastRun = lastRun.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func GetAllTasks() ([]core.ScheduledTask, error) {
	rows, err := DB.Query(`SELECT id, server_id, name, action, payload, cron_expression, one_shot, last_run FROM scheduled_tasks`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []core.ScheduledTask
	for rows.Next() {
		var t core.ScheduledTask
		var lastRun sql.NullTime
		if err := rows.Scan(&t.ID, &t.ServerID, &t.Name, &t.Action, &t.Payload, &t.CronExpression, &t.OneShot, &lastRun); err != nil {
			return nil, err
		}
		if lastRun.Valid {
			t.LastRun = lastRun.Time
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func DeleteTask(id string) error {
	_, err := DB.Exec(`DELETE FROM scheduled_tasks WHERE id = ?`, id)
	return err
}

func UpdateLastRun(id string, lastRun time.Time) error {
	_, err := DB.Exec(`UPDATE scheduled_tasks SET last_run = ? WHERE id = ?`, lastRun, id)
	return err
}
