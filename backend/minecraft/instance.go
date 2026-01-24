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

	Output chan string
}

func NewInstance(id string, runDir, jarName string) *Instance {
	return &Instance{
		ID:       id,
		RunDir:   runDir,
		JarName:  jarName,
		JavaArgs: []string{"-Xmx1G", "-Xms1G"},
		status:   core.StatusStopped,
		Output:   make(chan string, 100),
	}
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
	select {
	case i.Output <- msg:
	default:
	}
}
