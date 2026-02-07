package cmd

import (
	"testing"

	"github.com/linanwx/nagobot/config"
)

func TestIsOpenRouterKimiModel(t *testing.T) {
	tests := []struct {
		name string
		cfg  *config.Config
		want bool
	}{
		{
			name: "nil config",
			cfg:  nil,
			want: false,
		},
		{
			name: "openrouter kimi",
			cfg: &config.Config{
				Agents: config.AgentsConfig{
					Defaults: config.AgentDefaults{
						Provider:  "openrouter",
						ModelType: "moonshotai/kimi-k2.5",
					},
				},
			},
			want: true,
		},
		{
			name: "openrouter kimi uppercase",
			cfg: &config.Config{
				Agents: config.AgentsConfig{
					Defaults: config.AgentDefaults{
						Provider:  "OpenRouter",
						ModelType: "MOONSHOTAI/KIMI-K2.5",
					},
				},
			},
			want: true,
		},
		{
			name: "openrouter non kimi",
			cfg: &config.Config{
				Agents: config.AgentsConfig{
					Defaults: config.AgentDefaults{
						Provider:  "openrouter",
						ModelType: "deepseek-chat",
					},
				},
			},
			want: false,
		},
		{
			name: "non openrouter kimi",
			cfg: &config.Config{
				Agents: config.AgentsConfig{
					Defaults: config.AgentDefaults{
						Provider:  "anthropic",
						ModelType: "moonshotai/kimi-k2.5",
					},
				},
			},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := isOpenRouterKimiModel(tc.cfg)
			if got != tc.want {
				t.Fatalf("isOpenRouterKimiModel() = %v, want %v", got, tc.want)
			}
		})
	}
}
