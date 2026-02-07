package thread

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/linanwx/nagobot/logger"
	"github.com/linanwx/nagobot/provider"
	"github.com/linanwx/nagobot/tools"
)

// Runner is a generic agent loop executor.
type Runner struct {
	provider provider.Provider
	tools    *tools.Registry
}

// NewRunner creates a new Runner.
func NewRunner(p provider.Provider, t *tools.Registry) *Runner {
	return &Runner{
		provider: p,
		tools:    t,
	}
}

// RunWithMessages executes the agent loop with pre-built messages.
func (r *Runner) RunWithMessages(ctx context.Context, messages []provider.Message) (string, error) {
	toolDefs := r.tools.Defs()

	for {
		resp, err := r.provider.Chat(ctx, &provider.Request{
			Messages: messages,
			Tools:    toolDefs,
		})
		if err != nil {
			return "", fmt.Errorf("provider error: %w", err)
		}

		if !resp.HasToolCalls() {
			return resp.Content, nil
		}

		messages = append(messages, provider.AssistantMessageWithTools(resp.Content, resp.ReasoningContent, resp.ToolCalls))

		for _, tc := range resp.ToolCalls {
			result := r.tools.Run(ctx, tc.Function.Name, json.RawMessage(tc.Function.Arguments))
			if strings.HasPrefix(result, "Error:") {
				logger.Error("tool error", "tool", tc.Function.Name, "err", result)
			}
			messages = append(messages, provider.ToolResultMessage(tc.ID, tc.Function.Name, result))
		}
	}
}
