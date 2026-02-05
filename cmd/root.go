// Package cmd provides CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time.
	Version = "0.1.0"
)

// rootCmd is the root command.
var rootCmd = &cobra.Command{
	Use:   "nagobot",
	Short: "nagobot - a lightweight AI assistant",
	Long: `nagobot is a lightweight, Go-based AI assistant that supports
multiple LLM providers with precise model support.

Supported providers:
  - openrouter (Kimi K2.5)
  - anthropic (Claude Sonnet 4, Claude Opus 4)

Get started with: nagobot onboard`,
	Version: Version,
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
}
