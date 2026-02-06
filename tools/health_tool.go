package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	healthsnap "github.com/linanwx/nagobot/internal/health"
	"github.com/linanwx/nagobot/provider"
)

// HealthTool reports runtime health info for the current process.
type HealthTool struct{}

// Def returns the tool definition.
func (t *HealthTool) Def() provider.ToolDef {
	return provider.ToolDef{
		Type: "function",
		Function: provider.FunctionDef{
			Name:        "health",
			Description: "Get runtime health information (memory, goroutines, runtime metadata) for the current nagobot process.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"format": map[string]any{
						"type":        "string",
						"description": "Optional output format: 'json' or 'text'. Defaults to 'json'.",
						"enum":        []string{"json", "text"},
					},
				},
			},
		},
	}
}

type healthArgs struct {
	Format string `json:"format,omitempty"`
}

// Run executes the tool.
func (t *HealthTool) Run(ctx context.Context, args json.RawMessage) string {
	var a healthArgs
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return fmt.Sprintf("Error: invalid arguments: %v", err)
		}
	}

	snapshot := healthsnap.Collect()
	if strings.EqualFold(a.Format, "text") {
		return healthsnap.FormatText(snapshot)
	}

	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error: failed to serialize health snapshot: %v", err)
	}
	return string(data)
}
