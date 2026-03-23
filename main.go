package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	bm "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

func main() {
	keyPath := os.Getenv("HOST_KEY_PATH")
	if keyPath == "" {
		keyPath = "/data/host_key"
	}

	s, err := wish.NewServer(
		wish.WithAddress("0.0.0.0:22"),
		wish.WithHostKeyPath(keyPath),
		wish.WithMiddleware(
			bm.MiddlewareWithProgramHandler(teaHandler, termenv.TrueColor),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not create server: %v\n", err)
		os.Exit(1)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("SSH server listening on :22")
	go func() {
		if err := s.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		}
	}()

	<-done
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}
}

func teaHandler(s ssh.Session) *tea.Program {
	pty, _, ok := s.Pty()
	renderer := bm.MakeRenderer(s)
	m := NewModel(renderer)
	if ok {
		m.windowWidth = pty.Window.Width
		m.windowHeight = pty.Window.Height
	}
	return tea.NewProgram(m,
		tea.WithAltScreen(),
		tea.WithInput(s),
		tea.WithOutput(s),
	)
}
