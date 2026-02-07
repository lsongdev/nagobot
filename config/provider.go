package config

import (
	"errors"
	"os"
)

// GetModelName returns the effective model name (modelName or modelType).
func (c *Config) GetModelName() string {
	if c.Agents.Defaults.ModelName != "" {
		return c.Agents.Defaults.ModelName
	}
	return c.Agents.Defaults.ModelType
}

// GetAPIKey returns the API key for the configured provider.
func (c *Config) GetAPIKey() (string, error) {
	switch c.Agents.Defaults.Provider {
	case "openrouter":
		if v := os.Getenv("OPENROUTER_API_KEY"); v != "" {
			return v, nil
		}
		if c.Providers.OpenRouter == nil || c.Providers.OpenRouter.APIKey == "" {
			return "", errors.New("openrouter API key not configured")
		}
		return c.Providers.OpenRouter.APIKey, nil
	case "anthropic":
		if v := os.Getenv("ANTHROPIC_API_KEY"); v != "" {
			return v, nil
		}
		if c.Providers.Anthropic == nil || c.Providers.Anthropic.APIKey == "" {
			return "", errors.New("anthropic API key not configured")
		}
		return c.Providers.Anthropic.APIKey, nil
	default:
		return "", errors.New("unknown provider: " + c.Agents.Defaults.Provider)
	}
}

// GetAPIBase returns the API base URL for the configured provider (env overrides config).
func (c *Config) GetAPIBase() string {
	switch c.Agents.Defaults.Provider {
	case "openrouter":
		if v := os.Getenv("OPENROUTER_API_BASE"); v != "" {
			return v
		}
		if c.Providers.OpenRouter != nil && c.Providers.OpenRouter.APIBase != "" {
			return c.Providers.OpenRouter.APIBase
		}
	case "anthropic":
		if v := os.Getenv("ANTHROPIC_API_BASE"); v != "" {
			return v
		}
		if c.Providers.Anthropic != nil && c.Providers.Anthropic.APIBase != "" {
			return c.Providers.Anthropic.APIBase
		}
	}
	return ""
}
