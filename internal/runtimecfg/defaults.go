package runtimecfg

import "time"

const (
	AgentDefaultMaxTokens         = 8192
	AgentDefaultTemperature       = 0.95
	AgentDefaultMaxToolIterations = 20
	AgentSessionMaxMessages       = 40
)

const (
	BusDefaultBufferSize = 100
)

const (
	SubagentDefaultMaxConcurrent = 5
)

const (
	SubagentDefaultTimeout   = 5 * time.Minute
	SubagentWaitPollInterval = 100 * time.Millisecond
)

const (
	CLIChannelMessageBufferSize      = 10
	CLIChannelStopWaitTimeout        = 500 * time.Millisecond
	TelegramChannelMessageBufferSize = 100
	TelegramUpdateTimeoutSeconds     = 30
	TelegramMaxMessageLength         = 4096
)

const (
	ServeHeartbeatTickInterval = 60 * time.Second
)

const (
	ToolExecDefaultTimeoutSeconds = 60
	ToolExecOutputMaxChars        = 50000
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
	WorkspaceMemoryDirName = "memory"
	WorkspaceSkillsDirName = "skills"
)

const (
	MemoryTurnsDirName          = "turns"
	MemoryIndexDirName          = "index"
	MemoryGlobalSummaryFileName = "MEMORY.md"
	MemorySkillFileName         = "memory.md"
	MemoryIndexFileExt          = ".jsonl"
)

const (
	MemoryExcerptChars          = 50
	MemoryDailyMaxTurns         = 50
	MemoryDailySummaryMaxChars  = 500
	MemoryGlobalMaxTurns        = 50
	MemoryGlobalSummaryMaxChars = 1000
	MemoryMaxKeywordsPerTurn    = 8
	MemoryMaxMarkersPerTurn     = 6
)
