package process

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type Server struct {
	Cmd       *exec.Cmd
	Stdin     io.WriteCloser
	ServerDir string
	JarName   string
}

func NewServer(serverDir, jarName string) *Server {
	return &Server{
		ServerDir: serverDir,
		JarName:   jarName,
	}
}

func (s *Server) Start() error {
	// Command : java -Xmx1024M -Xms1024M -jar server.jar nogui
	s.Cmd = exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", s.JarName, "nogui")
	s.Cmd.Dir = s.ServerDir // IMPORTANT : We execute the command in the data folder

	// Connect pipes
	stdout, err := s.Cmd.StdoutPipe()
	if err != nil {
		return err
	}
	
	// We also capture the Stderr (Java errors)
	s.Cmd.Stderr = s.Cmd.Stdout

	stdin, err := s.Cmd.StdinPipe()
	if err != nil {
		return err
	}
	s.Stdin = stdin

	// Start the process
	if err := s.Cmd.Start(); err != nil {
		return err
	}

	fmt.Println("--- MINECRAFT SERVER STARTED ---")

	// Goroutine to read server logs in real time
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			// We display the server logs with a prefix
			fmt.Printf("[MC] %s\n", scanner.Text())
		}
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
