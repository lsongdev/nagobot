package provider

import (
	"fmt"
	"os"
	"strings"

	"github.com/linanwx/nagobot/config"
	"github.com/linanwx/nagobot/internal/runtimecfg"
)

// FactoryConfig stores provider-level credentials and endpoint settings.
type FactoryConfig struct {
	APIKey  string
	APIBase string
}

// Factory creates provider instances for the requested provider/model.
type Factory struct {
	configs          map[string]FactoryConfig
	defaultProv      string
	defaultModel     string
	defaultModelName string
	maxTokens        int
	temperature      float64
}

// NewFactory builds a provider factory from config.
func NewFactory(cfg *config.Config) (*Factory, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is nil")
	}

	defaultProv := strings.TrimSpace(cfg.Agents.Defaults.Provider)
	if defaultProv == "" {
		return nil, fmt.Errorf("default provider is required")
	}

	defaultModel := strings.TrimSpace(cfg.Agents.Defaults.ModelType)
	if defaultModel == "" {
		return nil, fmt.Errorf("default model type is required")
	}

	if err := ValidateProviderModelType(defaultProv, defaultModel); err != nil {
		return nil, err
	}

	maxTokens := cfg.Agents.Defaults.MaxTokens
	if maxTokens == 0 {
		maxTokens = runtimecfg.AgentDefaultMaxTokens
	}

	temperature := cfg.Agents.Defaults.Temperature
	if temperature == 0 {
		temperature = runtimecfg.AgentDefaultTemperature
	}

	f := &Factory{
		configs:          make(map[string]FactoryConfig),
		defaultProv:      defaultProv,
		defaultModel:     defaultModel,
		defaultModelName: cfg.GetModelName(),
		maxTokens:        maxTokens,
		temperature:      temperature,
	}

	for _, providerName := range SupportedProviders() {
		providerCfg := FactoryConfig{
			APIKey:  providerAPIKey(cfg, providerName),
			APIBase: providerAPIBase(cfg, providerName),
		}
		if providerCfg.APIKey != "" || providerName == defaultProv {
			f.configs[providerName] = providerCfg
		}
	}

	if conf, ok := f.configs[defaultProv]; !ok || strings.TrimSpace(conf.APIKey) == "" {
		return nil, fmt.Errorf("%s API key not configured", defaultProv)
	}

	return f, nil
}

// Create builds a provider instance for provider/model. Empty values fall back to defaults.
func (f *Factory) Create(providerName, modelType string) (Provider, error) {
	if f == nil {
		return nil, fmt.Errorf("provider factory is nil")
	}

	providerName = strings.TrimSpace(providerName)
	if providerName == "" {
		providerName = f.defaultProv
	}

	modelType = strings.TrimSpace(modelType)
	if modelType == "" {
		if providerName == f.defaultProv {
			modelType = f.defaultModel
		} else {
			models := SupportedModelsForProvider(providerName)
			if len(models) == 0 {
				return nil, fmt.Errorf("unknown provider: %s", providerName)
			}
			modelType = models[0]
		}
	}

	if err := ValidateProviderModelType(providerName, modelType); err != nil {
		return nil, err
	}

	provCfg, ok := f.configs[providerName]
	if !ok || strings.TrimSpace(provCfg.APIKey) == "" {
		return nil, fmt.Errorf("%s API key not configured", providerName)
	}

	modelName := modelType
	if providerName == f.defaultProv && modelType == f.defaultModel && strings.TrimSpace(f.defaultModelName) != "" {
		modelName = f.defaultModelName
	}

	switch providerName {
	case "openrouter":
		return NewOpenRouterProvider(provCfg.APIKey, provCfg.APIBase, modelType, modelName, f.maxTokens, f.temperature), nil
	case "anthropic":
		return NewAnthropicProvider(provCfg.APIKey, provCfg.APIBase, modelType, modelName, f.maxTokens, f.temperature), nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", providerName)
	}
}

func providerAPIKey(cfg *config.Config, providerName string) string {
	switch providerName {
	case "openrouter":
		if v := strings.TrimSpace(os.Getenv("OPENROUTER_API_KEY")); v != "" {
			return v
		}
		if cfg.Providers.OpenRouter != nil {
			return strings.TrimSpace(cfg.Providers.OpenRouter.APIKey)
		}
	case "anthropic":
		if v := strings.TrimSpace(os.Getenv("ANTHROPIC_API_KEY")); v != "" {
			return v
		}
		if cfg.Providers.Anthropic != nil {
			return strings.TrimSpace(cfg.Providers.Anthropic.APIKey)
		}
	}
	return ""
}

func providerAPIBase(cfg *config.Config, providerName string) string {
	switch providerName {
	case "openrouter":
		if v := strings.TrimSpace(os.Getenv("OPENROUTER_API_BASE")); v != "" {
			return v
		}
		if cfg.Providers.OpenRouter != nil {
			return strings.TrimSpace(cfg.Providers.OpenRouter.APIBase)
		}
	case "anthropic":
		if v := strings.TrimSpace(os.Getenv("ANTHROPIC_API_BASE")); v != "" {
			return v
		}
		if cfg.Providers.Anthropic != nil {
			return strings.TrimSpace(cfg.Providers.Anthropic.APIBase)
		}
	}
	return ""
}
