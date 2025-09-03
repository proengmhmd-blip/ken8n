package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea/v2"
	opencode "github.com/kenkaiii/ken8n-coder-sdk-go"
	"github.com/kenkaiii/ken8n-coder-sdk-go/option"
	"github.com/kenkaiii/ken8n-coder/internal/api"
	"github.com/kenkaiii/ken8n-coder/internal/app"
	"github.com/kenkaiii/ken8n-coder/internal/clipboard"
	"github.com/kenkaiii/ken8n-coder/internal/tui"
	"github.com/kenkaiii/ken8n-coder/internal/util"
	flag "github.com/spf13/pflag"
)

var Version = "1.1.5"

func main() {
	version := Version
	if version != "dev" && !strings.HasPrefix(Version, "v") {
		version = "v" + Version
	}

	var model *string = flag.String("model", "", "model to begin with")
	var prompt *string = flag.String("prompt", "", "prompt to begin with")
	var agent *string = flag.String("agent", "", "agent to begin with")
	var sessionID *string = flag.String("session", "", "session ID")
	flag.Parse()

	url := os.Getenv("KEN8N_CODER_SERVER")

	appInfoStr := os.Getenv("KEN8N_CODER_APP_INFO")
	var appInfo opencode.App
	err := json.Unmarshal([]byte(appInfoStr), &appInfo)
	if err != nil {
		slog.Error("Failed to unmarshal app info", "error", err)
		os.Exit(1)
	}

	stat, err := os.Stdin.Stat()
	if err != nil {
		slog.Error("Failed to stat stdin", "error", err)
		os.Exit(1)
	}

	// Check if there's data piped to stdin
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		stdin, err := io.ReadAll(os.Stdin)
		if err != nil {
			slog.Error("Failed to read stdin", "error", err)
			os.Exit(1)
		}
		stdinContent := strings.TrimSpace(string(stdin))
		if stdinContent != "" {
			if prompt == nil || *prompt == "" {
				prompt = &stdinContent
			} else {
				combined := *prompt + "\n" + stdinContent
				prompt = &combined
			}
		}
	}

	httpClient := opencode.NewClient(
		option.WithBaseURL(url),
	)

	// Fetch agents from the /agent endpoint
	agentsPtr, err := httpClient.App.Agents(context.Background())
	if err != nil {
		slog.Error("Failed to fetch agents", "error", err)
		os.Exit(1)
	}
	if agentsPtr == nil {
		slog.Error("No agents returned from server")
		os.Exit(1)
	}
	agents := *agentsPtr

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	apiHandler := util.NewAPILogHandler(ctx, httpClient, "tui", slog.LevelDebug)
	logger := slog.New(apiHandler)
	slog.SetDefault(logger)

	slog.Debug("TUI launched", "app", appInfoStr, "agents_count", len(agents), "url", url)

	go func() {
		err = clipboard.Init()
		if err != nil {
			slog.Error("Failed to initialize clipboard", "error", err)
		}
	}()

	// Create main context for the application
	app_, err := app.New(ctx, version, appInfo, agents, httpClient, model, prompt, agent, sessionID)
	if err != nil {
		panic(err)
	}

	tuiModel := tui.NewModel(app_).(*tui.Model)
	program := tea.NewProgram(
		tuiModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		stream := httpClient.Event.ListStreaming(ctx)
		for stream.Next() {
			evt := stream.Current().AsUnion()
			if _, ok := evt.(opencode.EventListResponseEventStorageWrite); ok {
				continue
			}
			program.Send(evt)
		}
		if err := stream.Err(); err != nil {
			slog.Error("Error streaming events", "error", err)
			program.Send(err)
		}
	}()

	go api.Start(ctx, program, httpClient)

	// Handle signals in a separate goroutine
	go func() {
		sig := <-sigChan
		slog.Info("Received signal, shutting down gracefully", "signal", sig)
		tuiModel.Cleanup()
		program.Quit()
	}()

	// Run the TUI
	result, err := program.Run()
	if err != nil {
		slog.Error("TUI error", "error", err)
	}

	tuiModel.Cleanup()
	slog.Info("TUI exited", "result", result)
}
