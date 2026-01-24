package process

import (
	"bufio"
	"io"
	"os/exec"
)

type Server struct {
	ID        int
	Port      int
	Cmd       *exec.Cmd
	Stdin     io.WriteCloser
	ServerDir string
	JarName   string
	Output    chan string
	IsRunning bool
}

func NewServer(serverDir, jarName string) *Server {
	return &Server{
		ServerDir: serverDir,
		JarName:   jarName,
		Output:    make(chan string, 100),
		IsRunning: false,
	}
}

func (s *Server) Start() error {
	if s.IsRunning {
		return nil
	}

	s.Cmd = exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", s.JarName, "nogui")
	s.Cmd.Dir = s.ServerDir

	stdout, err := s.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	s.Cmd.Stderr = s.Cmd.Stdout

	stdin, err := s.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	s.Stdin = stdin

	if err := s.Cmd.Start(); err != nil {
		return err
	}

	s.IsRunning = true
	s.broadcast("--- MINECRAFT SERVER STARTED ---")

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			text := scanner.Text()
			s.broadcast(text)
		}

		s.IsRunning = false
		s.broadcast("--- MINECRAFT SERVER STOPPED ---")
		s.Cmd.Wait()
	}()

	return nil
}

// WriteCommand sends a command to the Minecraft server
func (s *Server) WriteCommand(cmd string) error {
	_, err := io.WriteString(s.Stdin, cmd+"\n")
	return err
}

func (s *Server) Wait() error {
	return s.Cmd.Wait()
}

func (s *Server) broadcast(msg string) {
	select {
	case s.Output <- msg:
	default:
	}
}
