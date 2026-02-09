package config

import (
	"path/filepath"

	"github.com/linanwx/nagobot/internal/runtimecfg"
)

// DefaultConfig returns a config with sensible defaults.
func DefaultConfig() *Config {
	logDefaults := defaultLoggingConfig()
	return &Config{
		Thread: ThreadConfig{
			Provider:            runtimecfg.ThreadDefaultProvider,
			ModelType:           runtimecfg.ThreadDefaultModelType,
			MaxTokens:           runtimecfg.ThreadDefaultMaxTokens,
			Temperature:         runtimecfg.ThreadDefaultTemperature,
			ContextWindowTokens: runtimecfg.ThreadDefaultContextWindowTokens,
			ContextWarnRatio:    runtimecfg.ThreadDefaultContextWarnRatio,
		},
		Providers: ProvidersConfig{
			DeepSeek: &ProviderConfig{
				APIKey: "",
			},
		},
		Channels: &ChannelsConfig{
			AdminUserID: "",
			Telegram: &TelegramChannelConfig{
				Token:      "",
				AllowedIDs: []int64{},
			},
			Web: &WebChannelConfig{
				Addr: runtimecfg.WebChannelDefaultAddr,
			},
		},
		Logging: logDefaults,
	}
}

func defaultLoggingConfig() LoggingConfig {
	dir, err := ConfigDir()
	if err != nil {
		dir = ""
	}
	logFile := filepath.Join(dir, "logs", "nagobot.log")
	enabled := true
	return LoggingConfig{
		Enabled: &enabled,
		Level:   "info",
		Stdout:  true,
		File:    logFile,
	}
}

func (c *Config) applyDefaults() {
	if c.Thread.Provider == "" {
		c.Thread.Provider = runtimecfg.ThreadDefaultProvider
	}
	if c.Thread.ModelType == "" {
		c.Thread.ModelType = runtimecfg.ThreadDefaultModelType
	}
	if c.Thread.MaxTokens <= 0 {
		c.Thread.MaxTokens = runtimecfg.ThreadDefaultMaxTokens
	}
	if c.Thread.Temperature == 0 {
		c.Thread.Temperature = runtimecfg.ThreadDefaultTemperature
	}
	if c.Thread.ContextWindowTokens <= 0 {
		c.Thread.ContextWindowTokens = runtimecfg.ThreadDefaultContextWindowTokens
	}
	if c.Thread.ContextWarnRatio <= 0 || c.Thread.ContextWarnRatio >= 1 {
		c.Thread.ContextWarnRatio = runtimecfg.ThreadDefaultContextWarnRatio
	}

	if c.Channels == nil {
		c.Channels = &ChannelsConfig{}
	}
	if c.Channels.Telegram == nil {
		c.Channels.Telegram = &TelegramChannelConfig{
			AllowedIDs: []int64{},
		}
	}
	if c.Channels.Telegram.AllowedIDs == nil {
		c.Channels.Telegram.AllowedIDs = []int64{}
	}
	if c.Channels.Web == nil {
		c.Channels.Web = &WebChannelConfig{}
	}
	if c.Channels.Web.Addr == "" {
		c.Channels.Web.Addr = runtimecfg.WebChannelDefaultAddr
	}

	def := defaultLoggingConfig()
	if c.Logging == (LoggingConfig{}) {
		c.Logging = def
		return
	}

	hasAny := c.Logging.Level != "" || c.Logging.File != "" || c.Logging.Stdout
	if c.Logging.Enabled == nil && hasAny {
		enabled := true
		c.Logging.Enabled = &enabled
	}
	if c.Logging.Level == "" {
		c.Logging.Level = def.Level
	}
	if c.Logging.File == "" {
		c.Logging.File = def.File
	}
	if !c.Logging.Stdout && c.Logging.File == "" {
		c.Logging.Stdout = def.Stdout
	}
	if c.Logging.Enabled == nil {
		c.Logging.Enabled = def.Enabled
	}
}
