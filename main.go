package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

const (
	defaultHostKeyPath          = "/data/host_key"
	defaultSSHAddress           = "0.0.0.0:22"
	defaultShutdownTimeout      = 30 * time.Second
	defaultContentWarmTimeout   = 1500 * time.Millisecond
	defaultContentRefresh       = 5 * time.Minute
	shutdownWaitForServeErr     = 2 * time.Second
	shutdownTimeoutEnvVar       = "SHUTDOWN_TIMEOUT_SECONDS"
	contentRefreshEnvVar        = "NOTION_REFRESH_SECONDS"
	defaultShutdownTimeoutValue = "30"
)

type appConfig struct {
	address         string
	hostKeyPath     string
	shutdownTimeout time.Duration
	refreshInterval time.Duration
}

func loadConfig() appConfig {
	addr := os.Getenv("APP_ADDR")
	if addr == "" {
		addr = defaultSSHAddress
	}

	hostKeyPath := os.Getenv("HOST_KEY_PATH")
	if hostKeyPath == "" {
		hostKeyPath = defaultHostKeyPath
	}

	shutdownTimeout := defaultShutdownTimeout
	value := os.Getenv(shutdownTimeoutEnvVar)
	if value == "" {
		value = defaultShutdownTimeoutValue
	}

	refreshInterval := defaultContentRefresh
	refreshValue := os.Getenv(contentRefreshEnvVar)
	if refreshValue != "" {
		seconds, err := strconv.Atoi(refreshValue)
		if err != nil || seconds <= 0 {
			fmt.Fprintf(os.Stderr, "%s must be a positive integer, got %q. Falling back to %d seconds.\n", contentRefreshEnvVar, refreshValue, int(defaultContentRefresh/time.Second))
		} else {
			refreshInterval = time.Duration(seconds) * time.Second
		}
	}
	seconds, err := strconv.Atoi(value)
	if err != nil || seconds <= 0 {
		fmt.Fprintf(os.Stderr, "%s must be a positive integer, got %q. Falling back to %s seconds.\n", shutdownTimeoutEnvVar, value, defaultShutdownTimeoutValue)
	} else {
		shutdownTimeout = time.Duration(seconds) * time.Second
	}

	return appConfig{
		address:         addr,
		hostKeyPath:     hostKeyPath,
		shutdownTimeout: shutdownTimeout,
		refreshInterval: refreshInterval,
	}
}

func ensureHostKeyExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("host key path %q is not accessible: %w", path, err)
	}
	return nil
}

func main() {
	cfg := loadConfig()
	if err := ensureHostKeyExists(cfg.hostKeyPath); err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}

	s, err := wish.NewServer(
		wish.WithAddress(cfg.address),
		wish.WithHostKeyPath(cfg.hostKeyPath),
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

	appCtx, cancelApp := context.WithCancel(context.Background())
	defer cancelApp()

	appContentStore.WarmUp(defaultContentWarmTimeout)
	appContentStore.StartAutoRefresh(appCtx, cfg.refreshInterval)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- s.ListenAndServe()
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(done)

	fmt.Printf("SSH server listening on %s\n", cfg.address)

	select {
	case err := <-serveErr:
		if err != nil {
			fmt.Fprintf(os.Stderr, "server error: %v\n", err)
			os.Exit(1)
		}
		return
	case sig := <-done:
		fmt.Fprintf(os.Stderr, "received signal %s, shutting down\n", sig.String())
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.shutdownTimeout)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "shutdown error: %v\n", err)
	}

	select {
	case err := <-serveErr:
		if err != nil {
			fmt.Fprintf(os.Stderr, "server stopped with error during shutdown: %v\n", err)
		}
	case <-time.After(shutdownWaitForServeErr):
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
