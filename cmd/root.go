// Package cmd provides CLI commands.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/linanwx/nagobot/config"
	"github.com/linanwx/nagobot/logger"
	"github.com/linanwx/nagobot/provider"
	"github.com/spf13/cobra"
)

var (
	// Version is set at build time.
	Version = "0.1.0"

	logLevelOverride string
)

// rootCmd is the root command.
var rootCmd = &cobra.Command{
	Use:     "nagobot",
	Short:   "nagobot - a lightweight AI assistant",
	Long:    buildRootLong(),
	Version: Version,
}

func buildRootLong() string {
	var sb strings.Builder
	sb.WriteString("nagobot is a lightweight, Go-based AI assistant that supports\n")
	sb.WriteString("multiple LLM providers with precise model support.\n\n")
	sb.WriteString("Supported providers:\n")

	providers := provider.SupportedProviders()
	for _, name := range providers {
		models := provider.SupportedModelsForProvider(name)
		if len(models) > 0 {
			sb.WriteString(fmt.Sprintf("  - %s (%s)\n", name, strings.Join(models, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("  - %s\n", name))
		}
	}

	sb.WriteString("\nGet started with: nagobot onboard")
	return sb.String()
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&logLevelOverride, "log-level", "", "Override log level for this run (debug, info, warn, error)")
	rootCmd.PersistentPreRunE = applyRuntimeLogOverrides
}

func applyRuntimeLogOverrides(cmd *cobra.Command, args []string) error {
	if logLevelOverride == "" {
		return nil
	}

	level := strings.ToLower(strings.TrimSpace(logLevelOverride))
	switch level {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("invalid --log-level: %q (use debug, info, warn, error)", logLevelOverride)
	}

	cfg, err := config.Load()
	if err != nil {
		cfg = config.DefaultConfig()
	}
	cfg.Logging.Level = level

	configDir, _ := config.ConfigDir()
	logEnabled := true
	if cfg.Logging.Enabled != nil {
		logEnabled = *cfg.Logging.Enabled
	}

	logCfg := logger.Config{
		Enabled: logEnabled,
		Level:   cfg.Logging.Level,
		Stdout:  cfg.Logging.Stdout,
		File:    cfg.Logging.File,
	}

	if err := logger.Init(logCfg, configDir); err != nil {
		return fmt.Errorf("logger init error: %w", err)
	}
	return nil
}
