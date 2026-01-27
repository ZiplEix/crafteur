package minecraft

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sync"

	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ZiplEix/crafteur/core"
	"github.com/shirou/gopsutil/v3/process"
)

type WSMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type ServerStats struct {
	CpuUsage float64 `json:"cpu"`     // %
	RamUsage uint64  `json:"ram"`     // Octets
	RamMax   uint64  `json:"ram_max"` // Octets (Xmx)
}

type Instance struct {
	ID       string
	RunDir   string
	JarName  string
	JavaArgs []string

	cmd    *exec.Cmd
	stdin  io.WriteCloser
	status core.ServerStatus
	mu     sync.RWMutex

	subscribers []chan WSMessage
	subMu       sync.Mutex

	logs []string

	ConnectedPlayers map[string]bool
	playersMu        sync.RWMutex
}

func NewInstance(id string, runDir, jarName string) *Instance {
	return &Instance{
		ID:               id,
		RunDir:           runDir,
		JarName:          jarName,
		JavaArgs:         []string{"-Xmx1G", "-Xms1G"},
		status:           core.StatusStopped,
		subscribers:      make([]chan WSMessage, 0),
		logs:             make([]string, 0),
		ConnectedPlayers: make(map[string]bool),
	}
}

func (i *Instance) Subscribe() chan WSMessage {
	i.subMu.Lock()
	defer i.subMu.Unlock()
	ch := make(chan WSMessage, 100)
	i.subscribers = append(i.subscribers, ch)
	return ch
}

func (i *Instance) Unsubscribe(ch chan WSMessage) {
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
	i.status = status
	i.mu.Unlock()

	i.broadcast(WSMessage{Type: "status", Data: string(status)})
}

// Helper to check if player is online
func (i *Instance) IsPlayerOnline(name string) bool {
	i.playersMu.RLock()
	defer i.playersMu.RUnlock()
	return i.ConnectedPlayers[name]
}

func (i *Instance) SetRAM(mb int) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// Create formatted strings, e.g., "-Xmx1024M"
	xmx := fmt.Sprintf("-Xmx%dM", mb)
	xms := fmt.Sprintf("-Xms%dM", mb)

	// Update JavaArgs
	// We replace the default args or update existing ones
	// For simplicity, we'll assuming we control the args and just overwrite the memory ones
	// But to be safe, let's filter out old memory args and append new ones

	newArgs := make([]string, 0, len(i.JavaArgs)+2)
	for _, arg := range i.JavaArgs {
		if !strings.HasPrefix(arg, "-Xmx") && !strings.HasPrefix(arg, "-Xms") {
			newArgs = append(newArgs, arg)
		}
	}
	newArgs = append(newArgs, xmx, xms)
	i.JavaArgs = newArgs
}

var (
	joinRegex  = regexp.MustCompile(`]: (\w+) joined the game`)
	leaveRegex = regexp.MustCompile(`]: (\w+) left the game`)
)

func (i *Instance) Start() error {
	i.mu.Lock()
	if i.status != core.StatusStopped {
		i.mu.Unlock()
		return fmt.Errorf("server is already running")
	}
	i.status = core.StatusStarting
	i.mu.Unlock()

	i.broadcast(WSMessage{Type: "status", Data: string(core.StatusStarting)})

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
	i.broadcastLog("--- PROCESS START ---")

	go i.monitorProcess(stdout)
	go i.startMonitoring()

	return nil
}

func (i *Instance) monitorProcess(stdout io.Reader) {
	// Reset players on start
	i.playersMu.Lock()
	i.ConnectedPlayers = make(map[string]bool)
	i.playersMu.Unlock()

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		text := scanner.Text()
		i.broadcastLog(text)

		// Parse Join
		if matches := joinRegex.FindStringSubmatch(text); len(matches) > 1 {
			player := matches[1]
			i.playersMu.Lock()
			i.ConnectedPlayers[player] = true
			i.playersMu.Unlock()
		}

		// Parse Leave
		if matches := leaveRegex.FindStringSubmatch(text); len(matches) > 1 {
			player := matches[1]
			i.playersMu.Lock()
			delete(i.ConnectedPlayers, player)
			i.playersMu.Unlock()
		}
	}

	if err := i.cmd.Wait(); err != nil {
		i.broadcastLog(fmt.Sprintf("--- CRASH/STOP ERROR: %v ---", err))
	} else {
		i.broadcastLog("--- PROCESS STOPPED GRACEFULLY ---")
	}

	// Clear players on stop
	i.playersMu.Lock()
	i.ConnectedPlayers = make(map[string]bool)
	i.playersMu.Unlock()

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

func (i *Instance) broadcastLog(msg string) {
	i.subMu.Lock()
	defer i.subMu.Unlock()

	i.logs = append(i.logs, msg)
	if len(i.logs) > 100 {
		i.logs = i.logs[1:]
	}

	message := WSMessage{Type: "log", Data: msg}

	for _, ch := range i.subscribers {
		select {
		case ch <- message:
		default:
		}
	}
}

func (i *Instance) broadcast(msg WSMessage) {
	i.subMu.Lock()
	defer i.subMu.Unlock()

	for _, ch := range i.subscribers {
		select {
		case ch <- msg:
		default:
		}
	}
}

func (i *Instance) startMonitoring() {
	i.mu.Lock()
	if i.cmd == nil || i.cmd.Process == nil {
		i.mu.Unlock()
		return
	}
	pid := int32(i.cmd.Process.Pid)
	maxRam := i.parseMaxRam()
	i.mu.Unlock()

	proc, err := process.NewProcess(pid)
	if err != nil {
		return
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		if i.GetStatus() != core.StatusRunning {
			return
		}

		select {
		case <-ticker.C:
			exists, _ := process.PidExists(pid)
			if !exists {
				return
			}

			cpu, _ := proc.Percent(0)

			// Normalize CPU usage by core count to output 0-100% representation of system share
			// rather than total core accumulation (e.g. 800% on 8 cores)
			cpuNormalized := cpu / float64(runtime.NumCPU())

			mem, err := proc.MemoryInfo()
			var ramUsage uint64
			if err == nil {
				ramUsage = mem.RSS
			}

			stats := ServerStats{
				CpuUsage: cpuNormalized,
				RamUsage: ramUsage,
				RamMax:   maxRam,
			}

			i.broadcast(WSMessage{Type: "stats", Data: stats})
		}
	}
}

func (i *Instance) parseMaxRam() uint64 {
	for _, arg := range i.JavaArgs {
		if strings.HasPrefix(arg, "-Xmx") {
			return parseBytes(arg[4:])
		}
	}
	// Default 1G if not specified
	return 1024 * 1024 * 1024
}

func parseBytes(s string) uint64 {
	s = strings.ToUpper(s)
	var multiplier uint64 = 1

	if strings.HasSuffix(s, "G") {
		multiplier = 1024 * 1024 * 1024
		s = strings.TrimSuffix(s, "G")
	} else if strings.HasSuffix(s, "M") {
		multiplier = 1024 * 1024
		s = strings.TrimSuffix(s, "M")
	} else if strings.HasSuffix(s, "K") {
		multiplier = 1024
		s = strings.TrimSuffix(s, "K")
	}

	val, _ := strconv.ParseUint(s, 10, 64)
	return val * multiplier
}
