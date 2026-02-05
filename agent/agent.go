// Package agent implements the core agent loop and context management.
package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pinkplumcom/nagobot/config"
	"github.com/pinkplumcom/nagobot/provider"
	"github.com/pinkplumcom/nagobot/tools"
)

// Agent is the core agent that processes messages.
type Agent struct {
	cfg           *config.Config
	provider      provider.Provider
	tools         *tools.Registry
	workspace     string
	maxIterations int
}

// NewAgent creates a new agent.
func NewAgent(cfg *config.Config) (*Agent, error) {
	// Validate provider and model type
	if err := provider.ValidateProviderModelType(
		cfg.Agents.Defaults.Provider,
		cfg.Agents.Defaults.ModelType,
	); err != nil {
		return nil, err
	}

	// Get API key
	apiKey, err := cfg.GetAPIKey()
	if err != nil {
		return nil, err
	}

	// Get workspace
	workspace, err := cfg.WorkspacePath()
	if err != nil {
		return nil, err
	}

	// Create provider
	var p provider.Provider
	modelType := cfg.Agents.Defaults.ModelType
	modelName := cfg.GetModelName()
	maxTokens := cfg.Agents.Defaults.MaxTokens
	if maxTokens == 0 {
		maxTokens = 8192
	}
	temp := cfg.Agents.Defaults.Temperature
	if temp == 0 {
		temp = 0.7
	}

	switch cfg.Agents.Defaults.Provider {
	case "openrouter":
		p = provider.NewOpenRouterProvider(apiKey, modelType, modelName, maxTokens, temp)
	case "anthropic":
		p = provider.NewAnthropicProvider(apiKey, modelType, modelName, maxTokens, temp)
	default:
		return nil, errors.New("unknown provider: " + cfg.Agents.Defaults.Provider)
	}

	// Create tool registry
	registry := tools.NewRegistry()
	registry.RegisterDefaultTools(workspace)

	maxIter := cfg.Agents.Defaults.MaxToolIterations
	if maxIter == 0 {
		maxIter = 20
	}

	return &Agent{
		cfg:           cfg,
		provider:      p,
		tools:         registry,
		workspace:     workspace,
		maxIterations: maxIter,
	}, nil
}

// Run processes a user message and returns the assistant's response.
func (a *Agent) Run(ctx context.Context, userMessage string) (string, error) {
	// Build context
	messages := a.buildContext(userMessage)

	// Get tool definitions
	toolDefs := a.tools.Defs()

	// Agent loop
	for i := 0; i < a.maxIterations; i++ {
		// Call provider
		resp, err := a.provider.Chat(ctx, &provider.Request{
			Messages: messages,
			Tools:    toolDefs,
		})
		if err != nil {
			return "", fmt.Errorf("provider error: %w", err)
		}

		// No tool calls = done
		if !resp.HasToolCalls() {
			return resp.Content, nil
		}

		// Add assistant message with tool calls
		messages = append(messages, provider.AssistantMessageWithTools(resp.Content, resp.ToolCalls))

		// Execute tool calls
		for _, tc := range resp.ToolCalls {
			result := a.tools.Run(ctx, tc.Function.Name, tc.Arguments)
			messages = append(messages, provider.ToolResultMessage(tc.ID, tc.Function.Name, result))
		}
	}

	return "", errors.New("max iterations exceeded")
}

// buildContext builds the initial message context.
func (a *Agent) buildContext(userMessage string) []provider.Message {
	messages := []provider.Message{}

	// System prompt
	systemPrompt := a.buildSystemPrompt()
	messages = append(messages, provider.SystemMessage(systemPrompt))

	// User message
	messages = append(messages, provider.UserMessage(userMessage))

	return messages
}

// buildSystemPrompt builds the system prompt from workspace files.
func (a *Agent) buildSystemPrompt() string {
	var parts []string

	// Identity
	parts = append(parts, fmt.Sprintf(`# nagobot

You are nagobot, a helpful AI assistant. You have access to tools that allow you to:
- Read, write, and edit files
- List directory contents

## Current Time
%s

## Workspace
Your workspace is at: %s

All file operations should be relative to this workspace unless an absolute path is given.
`, time.Now().Format("2006-01-02 15:04 (Monday)"), a.workspace))

	// Bootstrap files
	bootstrapFiles := []string{"AGENTS.md", "SOUL.md", "USER.md", "IDENTITY.md"}
	for _, name := range bootstrapFiles {
		path := filepath.Join(a.workspace, name)
		content, err := os.ReadFile(path)
		if err == nil && len(content) > 0 {
			parts = append(parts, fmt.Sprintf("## %s\n\n%s", name, strings.TrimSpace(string(content))))
		}
	}

	// Memory
	memoryPath := filepath.Join(a.workspace, "memory", "MEMORY.md")
	memoryContent, err := os.ReadFile(memoryPath)
	if err == nil && len(memoryContent) > 0 {
		parts = append(parts, fmt.Sprintf("## Long-term Memory\n\n%s", strings.TrimSpace(string(memoryContent))))
	}

	// Today's notes
	todayFile := time.Now().Format("2006-01-02") + ".md"
	todayPath := filepath.Join(a.workspace, "memory", todayFile)
	todayContent, err := os.ReadFile(todayPath)
	if err == nil && len(todayContent) > 0 {
		parts = append(parts, fmt.Sprintf("## Today's Notes\n\n%s", strings.TrimSpace(string(todayContent))))
	}

	// Available tools
	toolNames := a.tools.Names()
	parts = append(parts, fmt.Sprintf("## Available Tools\n\n%s", strings.Join(toolNames, ", ")))

	return strings.Join(parts, "\n\n---\n\n")
}

// ============================================================================
// Session management (simple in-memory + file persistence)
// ============================================================================

// Session represents a conversation session.
type Session struct {
	Key       string             `json:"key"`
	Messages  []provider.Message `json:"messages"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// SessionManager manages conversation sessions.
type SessionManager struct {
	sessionsDir string
	cache       map[string]*Session
}

// NewSessionManager creates a new session manager.
func NewSessionManager(configDir string) (*SessionManager, error) {
	sessionsDir := filepath.Join(configDir, "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return nil, err
	}
	return &SessionManager{
		sessionsDir: sessionsDir,
		cache:       make(map[string]*Session),
	}, nil
}

// Get returns a session by key, creating one if it doesn't exist.
func (m *SessionManager) Get(key string) (*Session, error) {
	// Check cache
	if s, ok := m.cache[key]; ok {
		return s, nil
	}

	// Try to load from disk
	path := m.sessionPath(key)
	data, err := os.ReadFile(path)
	if err == nil {
		var s Session
		if err := json.Unmarshal(data, &s); err == nil {
			m.cache[key] = &s
			return &s, nil
		}
	}

	// Create new session
	s := &Session{
		Key:       key,
		Messages:  []provider.Message{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.cache[key] = s
	return s, nil
}

// Save saves a session to disk.
func (m *SessionManager) Save(s *Session) error {
	s.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.sessionPath(s.Key), data, 0644)
}

// sessionPath returns the file path for a session.
func (m *SessionManager) sessionPath(key string) string {
	// Sanitize key for filename
	safe := strings.ReplaceAll(key, "/", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	return filepath.Join(m.sessionsDir, safe+".json")
}

// ============================================================================
// Memory management
// ============================================================================

// Memory manages long-term and daily memory.
type Memory struct {
	memoryDir string
}

// NewMemory creates a new memory manager.
func NewMemory(workspace string) (*Memory, error) {
	memoryDir := filepath.Join(workspace, "memory")
	if err := os.MkdirAll(memoryDir, 0755); err != nil {
		return nil, err
	}
	return &Memory{memoryDir: memoryDir}, nil
}

// ReadLongTerm reads the long-term memory.
func (m *Memory) ReadLongTerm() (string, error) {
	path := filepath.Join(m.memoryDir, "MEMORY.md")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(content), nil
}

// WriteLongTerm writes to the long-term memory.
func (m *Memory) WriteLongTerm(content string) error {
	path := filepath.Join(m.memoryDir, "MEMORY.md")
	return os.WriteFile(path, []byte(content), 0644)
}

// ReadToday reads today's notes.
func (m *Memory) ReadToday() (string, error) {
	path := filepath.Join(m.memoryDir, time.Now().Format("2006-01-02")+".md")
	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}
	return string(content), nil
}

// AppendToday appends to today's notes.
func (m *Memory) AppendToday(content string) error {
	path := filepath.Join(m.memoryDir, time.Now().Format("2006-01-02")+".md")

	// Read existing content
	existing, _ := os.ReadFile(path)

	// If file doesn't exist, add header
	if len(existing) == 0 {
		header := fmt.Sprintf("# %s\n\n", time.Now().Format("2006-01-02"))
		existing = []byte(header)
	}

	// Append new content
	newContent := string(existing) + content + "\n"
	return os.WriteFile(path, []byte(newContent), 0644)
}
