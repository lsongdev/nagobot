package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/linanwx/nagobot/config"
	"github.com/linanwx/nagobot/provider"
)

func printModelRoutingNotice(cfg *config.Config) {
	if !isOpenRouterKimiModel(cfg) {
		return
	}

	fmt.Fprintln(
		os.Stderr,
		"[nagobot] Warning: OpenRouter + Kimi must use OpenRouter's official `moonshot` provider route. "+
			"Set provider allowlist to `moonshot` (or use a Moonshot-pinned OpenRouter alias), otherwise reasoning/tool support is not guaranteed and tool calls may fail.",
	)
}

func isOpenRouterKimiModel(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	if !strings.EqualFold(strings.TrimSpace(cfg.GetProvider()), "openrouter") {
		return false
	}

	modelType := strings.ToLower(strings.TrimSpace(cfg.GetModelType()))
	return provider.IsKimiModel(modelType)
}
