package services

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ZiplEix/crafteur/core"
	"github.com/ZiplEix/crafteur/database"
	"github.com/robfig/cron/v3"
)

type SchedulerService struct {
	cron          *cron.Cron
	serverService *ServerService
	entryMap      map[string]cron.EntryID // Maps TaskID to CronEntryID
	mu            sync.Mutex
}

func NewSchedulerService(serverService *ServerService) *SchedulerService {
	return &SchedulerService{
		cron:          cron.New(),
		serverService: serverService,
		entryMap:      make(map[string]cron.EntryID),
	}
}

func (s *SchedulerService) Start() {
	s.cron.Start()
}

func (s *SchedulerService) Stop() {
	s.cron.Stop()
}

func (s *SchedulerService) LoadTasks() error {
	tasks, err := database.GetAllTasks()
	if err != nil {
		return err
	}

	fmt.Printf("Chargement de %d tâches planifiées...\n", len(tasks))
	for _, task := range tasks {
		if err := s.ScheduleTask(&task); err != nil {
			fmt.Printf("Erreur chargement tâche %s: %v\n", task.ID, err)
		}
	}
	return nil
}

func (s *SchedulerService) ScheduleTask(task *core.ScheduledTask) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Remove existing if any (update case)
	if entryID, exists := s.entryMap[task.ID]; exists {
		s.cron.Remove(entryID)
	}

	entryID, err := s.cron.AddFunc(task.CronExpression, func() {
		s.executeTask(task)
	})

	if err != nil {
		return err
	}

	s.entryMap[task.ID] = entryID
	return nil
}

func (s *SchedulerService) UnscheduleTask(taskID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if entryID, exists := s.entryMap[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.entryMap, taskID)
	}
}

func (s *SchedulerService) executeTask(task *core.ScheduledTask) {
	fmt.Printf("Exécution tâche planifiée: %s (ID: %s) - Action: %s\n", task.Name, task.ID, task.Action)

	var err error
	switch task.Action {
	case "start":
		err = s.serverService.StartServer(task.ServerID)
	case "stop":
		err = s.serverService.StopServer(task.ServerID)
		// Wait a bit if we want to restart later? No, restart is a separate action or handled by script.
	case "restart":
		// Stop, wait, start
		err = s.serverService.StopServer(task.ServerID)
		if err == nil {
			time.Sleep(5 * time.Second)
			err = s.serverService.StartServer(task.ServerID)
		}
	case "command":
		commands := strings.Split(task.Payload, "\n")
		for _, cmd := range commands {
			if strings.TrimSpace(cmd) != "" {
				if e := s.serverService.SendCommand(task.ServerID, cmd); e != nil {
					err = e // Keep last error
				}
			}
		}
	default:
		err = fmt.Errorf("unknown action: %s", task.Action)
	}

	if err != nil {
		fmt.Printf("Erreur exécution tâche %s: %v\n", task.ID, err)
	}

	// Update LastRun
	database.UpdateLastRun(task.ID, time.Now())

	// Handle OneShot
	if task.OneShot {
		s.UnscheduleTask(task.ID)
		database.DeleteTask(task.ID)
	}
}

func (s *SchedulerService) CreateTask(task *core.ScheduledTask) error {
	if err := database.CreateTask(task); err != nil {
		return err
	}
	return s.ScheduleTask(task)
}

func (s *SchedulerService) DeleteTask(id string) error {
	if err := database.DeleteTask(id); err != nil {
		return err
	}
	s.UnscheduleTask(id)
	return nil
}

func (s *SchedulerService) GetTasksByServer(serverID string) ([]core.ScheduledTask, error) {
	tasks, err := database.GetTasksByServer(serverID)
	if err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Fill NextRun
	for i := range tasks {
		if entryID, exists := s.entryMap[tasks[i].ID]; exists {
			entry := s.cron.Entry(entryID)
			tasks[i].NextRun = entry.Next
		}
	}

	return tasks, nil
}
