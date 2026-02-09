package runtimecfg

import "time"

const (
	ThreadDefaultProvider            = "deepseek"
	ThreadDefaultModelType           = "deepseek-reasoner"
	ThreadDefaultMaxTokens           = 8192
	ThreadDefaultTemperature         = 0.95
	ThreadDefaultContextWindowTokens = 128000
	ThreadDefaultContextWarnRatio    = 0.8
)

const (
	CLIChannelMessageBufferSize      = 10
	CLIChannelStopWaitTimeout        = 500 * time.Millisecond
	TelegramChannelMessageBufferSize = 100
	TelegramUpdateTimeoutSeconds     = 30
	TelegramMaxMessageLength         = 4096
	WebChannelMessageBufferSize      = 100
	WebChannelDefaultAddr            = "127.0.0.1:8080"
	WebChannelShutdownTimeout        = 5 * time.Second
)

const (
	ToolExecDefaultTimeoutSeconds = 60
	ToolExecOutputMaxChars        = 50000
	ToolResultMaxChars            = 100000
)

const (
	ToolWebSearchDefaultMaxResults = 5
	ToolWebSearchHTTPTimeout       = 15 * time.Second
)

const (
	ToolWebFetchHTTPTimeout     = 30 * time.Second
	ToolWebFetchMaxReadBytes    = 500000
	ToolWebFetchMaxContentChars = 100000
)

const (
	ProviderSDKMaxRetries      = 2
	AnthropicFallbackMaxTokens = 1024
)

const (
	WorkspaceSkillsDirName   = "skills"
	WorkspaceSessionsDirName = "sessions"
)
