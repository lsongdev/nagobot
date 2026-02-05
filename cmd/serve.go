package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"path/filepath"
	"time"

	"github.com/linanwx/nagobot/agent"
	"github.com/linanwx/nagobot/channel"
	"github.com/linanwx/nagobot/config"
	"github.com/linanwx/nagobot/cron"
	"github.com/linanwx/nagobot/logger"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start nagobot as a service with channel integrations",
	Long: `Start nagobot as a long-running service that listens on multiple channels.

Supported channels:
  - cli: Interactive command line (default)
  - telegram: Telegram bot (requires TELEGRAM_BOT_TOKEN)

Examples:
  nagobot serve              # Start with CLI channel
  nagobot serve --telegram   # Start with Telegram bot
  nagobot serve --all        # Start all configured channels`,
	RunE: runServe,
}

var (
	serveTelegram bool
	serveAll      bool
	serveCLI      bool
)

func init() {
	serveCmd.Flags().BoolVar(&serveTelegram, "telegram", false, "Enable Telegram bot channel")
	serveCmd.Flags().BoolVar(&serveAll, "all", false, "Enable all configured channels")
	serveCmd.Flags().BoolVar(&serveCLI, "cli", true, "Enable CLI channel (default: true)")
	rootCmd.AddCommand(serveCmd)
}

func runServe(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create agent
	ag, err := agent.NewAgent(cfg)
	if err != nil {
		return fmt.Errorf("failed to create agent: %w", err)
	}
	defer ag.Close()

	// Create message handler
	handler := func(ctx context.Context, msg *channel.Message) (*channel.Response, error) {
		logger.Info("received message",
			"channel", msg.ChannelID,
			"user", msg.Username,
			"text", truncate(msg.Text, 50),
		)

		// Run agent
		response, err := ag.Run(ctx, msg.Text)
		if err != nil {
			logger.Error("agent error", "err", err)
			return &channel.Response{
				Text:    fmt.Sprintf("Error: %v", err),
				ReplyTo: msg.Metadata["chat_id"],
			}, nil
		}

		return &channel.Response{
			Text:    response,
			ReplyTo: msg.Metadata["chat_id"],
		}, nil
	}

	// Create router and manager
	router := channel.NewRouter(handler)
	manager := channel.NewManager(router)

	// Register channels
	if serveCLI || (!serveTelegram && !serveAll) {
		manager.Register(channel.NewCLIChannel(channel.CLIConfig{
			Prompt: "nagobot> ",
		}))
		logger.Info("CLI channel enabled")
	}

	if serveTelegram || serveAll {
		token := os.Getenv("TELEGRAM_BOT_TOKEN")
		if token == "" {
			// Try from config
			if cfg.Channels != nil && cfg.Channels.Telegram != nil {
				token = cfg.Channels.Telegram.Token
			}
		}

		if token == "" {
			logger.Warn("Telegram token not configured, skipping Telegram channel")
		} else {
			var allowedIDs []int64
			if cfg.Channels != nil && cfg.Channels.Telegram != nil {
				allowedIDs = cfg.Channels.Telegram.AllowedIDs
			}

			manager.Register(channel.NewTelegramChannel(channel.TelegramConfig{
				Token:      token,
				AllowedIDs: allowedIDs,
			}))
			logger.Info("Telegram channel enabled")
		}
	}

	// Setup context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		logger.Info("shutdown signal received")
		cancel()
	}()

	// Create and start cron service
	workspace, err := cfg.WorkspacePath()
	if err != nil {
		return fmt.Errorf("failed to get workspace: %w", err)
	}
	cronPath := filepath.Join(workspace, "cron.yaml")
	cronSvc, err := cron.NewService(cronPath)
	if err != nil {
		logger.Warn("cron service not started", "err", err)
	}
	if cronSvc != nil {
		cronSvc.Start(func(ctx context.Context, job cron.Job) error {
			logger.Info("cron job triggered", "id", job.ID, "name", job.Name)
			_, err := ag.Run(ctx, job.Payload.Message)
			return err
		})
		if err := cronSvc.StartWatching(); err != nil {
			logger.Warn("cron file watcher failed", "err", err)
		}
	}

	// Create and start heartbeat monitor
	heartbeat := cron.NewHeartbeat(cron.HeartbeatConfig{
		Interval: 30 * time.Second,
		OnUnhealthy: func(health *cron.SystemHealth) {
			logger.Error("system unhealthy",
				"status", health.Status,
				"goroutines", health.Goroutines,
				"memoryMB", health.Memory.Alloc/1024/1024,
			)
		},
	})

	// Register health checks
	heartbeat.Register("agent", cron.PingCheck())
	heartbeat.Register("memory", cron.MemoryCheck(500*1024*1024)) // 500MB limit
	heartbeat.Register("goroutines", cron.GoroutineCheck(1000))   // 1000 goroutine limit

	if err := heartbeat.Start(ctx); err != nil {
		return fmt.Errorf("failed to start heartbeat: %w", err)
	}

	// Start all channels
	if err := manager.StartAll(ctx); err != nil {
		return fmt.Errorf("failed to start channels: %w", err)
	}

	logger.Info("nagobot service started")
	fmt.Println("nagobot is running. Press Ctrl+C to stop.")

	// Wait for context cancellation
	<-ctx.Done()

	// Stop all components
	if err := manager.StopAll(); err != nil {
		logger.Error("error stopping channels", "err", err)
	}

	if cronSvc != nil {
		cronSvc.Stop()
	}

	if err := heartbeat.Stop(); err != nil {
		logger.Error("error stopping heartbeat", "err", err)
	}

	logger.Info("nagobot service stopped")
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
