package core

import "time"

type ScheduledTask struct {
	ID             string    `json:"id"`
	ServerID       string    `json:"server_id"`
	Name           string    `json:"name"`
	Action         string    `json:"action"`          // "start", "stop", "restart", "command"
	Payload        string    `json:"payload"`         // Command to execute
	CronExpression string    `json:"cron_expression"` // e.g., "0 10 * * *" or "@every 1h"
	OneShot        bool      `json:"one_shot"`        // If true, delete after execution
	LastRun        time.Time `json:"last_run"`
	NextRun        time.Time `json:"next_run"` // Not stored in DB
}
