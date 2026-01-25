package minecraft

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/ZiplEix/crafteur/core"
)

type Instance struct {
	ID       string
	RunDir   string
	JarName  string
	JavaArgs []string

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	status core.ServerStatus
	mu     sync.RWMutex

	subscribers []chan string
	subMu       sync.Mutex

	logs []string
}

func NewInstance(id string, runDir, jarName string) *Instance {
	return &Instance{
		ID:          id,
		RunDir:      runDir,
		JarName:     jarName,
		JavaArgs:    []string{"-Xmx1G", "-Xms1G"},
		status:      core.StatusStopped,
		subscribers: make([]chan string, 0),
		logs:        make([]string, 0),
	}
}

func (i *Instance) Subscribe() chan string {
	i.subMu.Lock()
	defer i.subMu.Unlock()
	ch := make(chan string, 100)
	i.subscribers = append(i.subscribers, ch)
	return ch
}

func (i *Instance) Unsubscribe(ch chan string) {
	i.subMu.Lock()
	defer i.subMu.Unlock()
	for idx, sub := range i.subscribers {
		if sub == ch {
			i.subscribers = append(i.subscribers[:idx], i.subscribers[idx+1:]...)
			close(ch)
			break
		}
	}
}

func (i *Instance) GetHistory() []string {
	i.subMu.Lock()
	defer i.subMu.Unlock()
	// Return a copy to be safe
	history := make([]string, len(i.logs))
	copy(history, i.logs)
	return history
}

func (i *Instance) GetStatus() core.ServerStatus {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.status
}

func (i *Instance) SetStatus(status core.ServerStatus) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.status = status
}

func (i *Instance) Start() error {
	i.mu.Lock()
	if i.status != core.StatusStopped {
		i.mu.Unlock()
		return fmt.Errorf("server is already running")
	}
	i.status = core.StatusStarting
	i.mu.Unlock()

	args := append(i.JavaArgs, "-jar", i.JarName, "nogui")
	i.cmd = exec.Command("java", args...)
	i.cmd.Dir = i.RunDir

	stdout, err := i.cmd.StdoutPipe()
	if err != nil {
		i.SetStatus(core.StatusStopped)
		return err
	}
	i.cmd.Stderr = i.cmd.Stdout

	stdin, err := i.cmd.StdinPipe()
	if err != nil {
		i.SetStatus(core.StatusStopped)
		return err
	}
	i.stdin = stdin

	if err := i.cmd.Start(); err != nil {
		i.SetStatus(core.StatusStopped)
		return err
	}

	i.SetStatus(core.StatusRunning)
	i.broadcast("--- PROCESS START ---")

	go i.monitorProcess(stdout)

	return nil
}

func (i *Instance) monitorProcess(stdout io.Reader) {
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		i.broadcast(text)
	}

	if err := i.cmd.Wait(); err != nil {
		i.broadcast(fmt.Sprintf("--- CRASH/STOP ERROR: %v ---", err))
	} else {
		i.broadcast("--- PROCESS STOPPED GRACEFULLY ---")
	}

	i.SetStatus(core.StatusStopped)
	i.cmd = nil
	i.stdin = nil
}

func (i *Instance) Stop() error {
	if i.GetStatus() == core.StatusStopped {
		return nil
	}

	err := i.SendCommand("stop")
	if err == nil {
		i.SetStatus(core.StatusStopping)
		return nil
	}

	if i.cmd != nil && i.cmd.Process != nil {
		return i.cmd.Process.Kill()
	}
	return fmt.Errorf("instance not found")
}

func (i *Instance) SendCommand(cmd string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if i.stdin == nil || i.status != core.StatusRunning {
		return fmt.Errorf("server stopped")
	}
	_, err := io.WriteString(i.stdin, cmd+"\n")

	return err
}

func (i *Instance) broadcast(msg string) {
	i.subMu.Lock()
	defer i.subMu.Unlock()

	i.logs = append(i.logs, msg)
	if len(i.logs) > 100 {
		i.logs = i.logs[1:]
	}

	for _, ch := range i.subscribers {
		select {
		case ch <- msg:
		default:
		}
	}
}
