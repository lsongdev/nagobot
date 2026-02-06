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
	"sync"
	"time"

	"github.com/linanwx/nagobot/bus"
	"github.com/linanwx/nagobot/config"
	"github.com/linanwx/nagobot/internal/runtimecfg"
	"github.com/linanwx/nagobot/logger"
	"github.com/linanwx/nagobot/provider"
	"github.com/linanwx/nagobot/skills"
	"github.com/linanwx/nagobot/tools"
)

// pendingSubagentResult holds a subagent completion or error for injection into the next turn.
type pendingSubagentResult struct {
	TaskID string
	Result string
	Error  string
}

// Agent is the core agent that processes messages.
type Agent struct {
	id            string
	cfg           *config.Config
	provider      provider.Provider
	tools         *tools.Registry
	skills        *skills.Registry
	agents        *AgentRegistry // Agent definitions from agents/ directory
	subagents     *bus.SubagentManager
	workspace     string
	maxIterations int
	runner        *Runner // Reusable runner for agent loop
	toolsDefaults tools.DefaultToolsConfig
	sessions      *SessionManager
	memory        *memoryStore

	// pendingResults keyed by session key to prevent cross-session leakage.
	// Empty key "" is used for stateless/CLI mode.
	pendingResults   map[string][]pendingSubagentResult
	pendingResultsMu sync.Mutex

	channelSender tools.ChannelSender // Set in serve mode for push delivery
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
	apiBase := cfg.GetAPIBase()

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
		maxTokens = runtimecfg.AgentDefaultMaxTokens
	}
	temp := cfg.Agents.Defaults.Temperature
	if temp == 0 {
		temp = runtimecfg.AgentDefaultTemperature
	}

	switch cfg.Agents.Defaults.Provider {
	case "openrouter":
		p = provider.NewOpenRouterProvider(apiKey, apiBase, modelType, modelName, maxTokens, temp)
	case "anthropic":
		p = provider.NewAnthropicProvider(apiKey, apiBase, modelType, modelName, maxTokens, temp)
	default:
		return nil, errors.New("unknown provider: " + cfg.Agents.Defaults.Provider)
	}

	// Generate agent ID
	agentID := fmt.Sprintf("agent-%d", time.Now().UnixNano())

	// Create tool registry
	toolRegistry := tools.NewRegistry()
	toolsDefaults := tools.DefaultToolsConfig{
		ExecTimeout:         cfg.Tools.Exec.Timeout,
		WebSearchMaxResults: cfg.Tools.Web.Search.MaxResults,
		RestrictToWorkspace: cfg.Tools.Exec.RestrictToWorkspace,
	}
	toolRegistry.RegisterDefaultTools(workspace, toolsDefaults)

	// Create skill registry
	skillRegistry := skills.NewRegistry()
	// Load custom skills from workspace
	skillsDir := filepath.Join(workspace, runtimecfg.WorkspaceSkillsDirName)
	if err := skillRegistry.LoadFromDirectory(skillsDir); err != nil {
		logger.Warn("failed to load skills", "dir", skillsDir, "err", err)
	}

	// Create agent registry (loads from agents/ directory)
	agentRegistry := NewAgentRegistry(workspace)

	maxIter := cfg.Agents.Defaults.MaxToolIterations
	if maxIter == 0 {
		maxIter = runtimecfg.AgentDefaultMaxToolIterations
	}

	// Create the runner (used by both main agent and subagents)
	runner := NewRunner(p, toolRegistry, maxIter)

	// Create session manager (non-fatal if it fails)
	configDir, _ := config.ConfigDir()
	sessions, sessErr := NewSessionManager(configDir)
	if sessErr != nil {
		logger.Warn("session manager unavailable", "err", sessErr)
	}

	agent := &Agent{
		id:             agentID,
		cfg:            cfg,
		provider:       p,
		tools:          toolRegistry,
		skills:         skillRegistry,
		agents:         agentRegistry,
		workspace:      workspace,
		maxIterations:  maxIter,
		runner:         runner,
		toolsDefaults:  toolsDefaults,
		sessions:       sessions,
		memory:         newMemoryStore(workspace),
		pendingResults: make(map[string][]pendingSubagentResult),
	}

	// Create subagent manager with a completion callback.
	// Push to channel if possible, otherwise queue for next turn.
	subagentMgr := bus.NewSubagentManager(5, agent.createSubagentRunner(), func(task *bus.SubagentTask) {
		var text string
		if task.Error != "" {
			text = fmt.Sprintf("[Subagent %s failed]: %s", task.ID, task.Error)
		} else {
			text = fmt.Sprintf("[Subagent %s completed]:\n%s", task.ID, task.Result)
		}

		// Attempt push delivery if origin and sender are available
		o := task.Origin
		if o.Channel != "" && o.ReplyTo != "" && agent.channelSender != nil {
			if err := agent.channelSender.SendTo(context.Background(), o.Channel, text, o.ReplyTo); err != nil {
				logger.Warn("subagent push delivery failed, queuing", "err", err)
			} else {
				return
			}
		}

		// Fallback: queue for next user message in this session
		agent.pendingResultsMu.Lock()
		agent.pendingResults[o.SessionKey] = append(agent.pendingResults[o.SessionKey], pendingSubagentResult{
			TaskID: task.ID,
			Result: task.Result,
			Error:  task.Error,
		})
		agent.pendingResultsMu.Unlock()
	})
	agent.subagents = subagentMgr

	// Register subagent tools (now that subagentMgr is ready)
	toolRegistry.Register(tools.NewSpawnAgentTool(subagentMgr, agentID))
	toolRegistry.Register(tools.NewCheckAgentTool(subagentMgr))

	// Register skill tool for progressive loading
	toolRegistry.Register(tools.NewUseSkillTool(skillRegistry))

	return agent, nil
}

// SetChannelSender registers the send_message tool backed by the given sender,
// and enables push delivery for async subagent results.
func (a *Agent) SetChannelSender(sender tools.ChannelSender) {
	a.channelSender = sender
	a.tools.Register(tools.NewSendMessageTool(sender))
}

// Close cleans up agent resources.
func (a *Agent) Close() {
}

// ID returns the agent's ID.
func (a *Agent) ID() string {
	return a.id
}

// createSubagentRunner creates a runner function for subagents.
// This runner reuses the same provider but creates a separate tool registry
// (without spawn_agent to prevent recursive spawning).
func (a *Agent) createSubagentRunner() bus.SubagentRunner {
	return func(ctx context.Context, task *bus.SubagentTask) (string, error) {
		// task.Type contains the agent name (from agents/ directory)
		agentName := task.Type

		// Create a tool registry for the subagent (without spawn tools to prevent recursion)
		subTools := tools.NewRegistry()
		subTools.RegisterDefaultTools(a.workspace, a.toolsDefaults)
		// Note: we don't register spawn_agent and check_agent to prevent infinite recursion

		// Create a runner for this subagent (same max iterations as main agent)
		subRunner := NewRunner(a.provider, subTools, a.maxIterations)

		// Build subagent system prompt from agents/ directory
		systemPrompt, err := a.agents.BuildPrompt(agentName, a.workspace, subTools.Names(), task.Task)
		if err != nil {
			return "", err
		}

		// Build the user message (task + optional context)
		userMessage := task.Task
		if task.Context != "" {
			userMessage = fmt.Sprintf("%s\n\nContext:\n%s", task.Task, task.Context)
		}

		// Execute using the runner
		return subRunner.Run(ctx, systemPrompt, userMessage)
	}
}

// drainPendingResults collects and clears pending subagent results for a specific session.
// Returns a formatted string to prepend to the user message, or empty if none.
func (a *Agent) drainPendingResults(sessionKey string) string {
	a.pendingResultsMu.Lock()
	results := a.pendingResults[sessionKey]
	if len(results) > 0 {
		delete(a.pendingResults, sessionKey)
	}
	a.pendingResultsMu.Unlock()

	if len(results) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("[Async subagent results since last turn]\n")
	for _, r := range results {
		if r.Error != "" {
			sb.WriteString(fmt.Sprintf("- %s failed: %s\n", r.TaskID, r.Error))
		} else {
			sb.WriteString(fmt.Sprintf("- %s completed: %s\n", r.TaskID, r.Result))
		}
	}
	return sb.String()
}

// Run processes a user message and returns the assistant's response (stateless).
func (a *Agent) Run(ctx context.Context, userMessage string) (string, error) {
	// Prepend any pending subagent results (stateless = empty session key)
	if prefix := a.drainPendingResults(""); prefix != "" {
		userMessage = prefix + "---\nUser message: " + userMessage
	}

	// Build system prompt
	systemPrompt := a.buildSystemPrompt()

	// Use the runner to execute the agent loop
	return a.runner.Run(ctx, systemPrompt, userMessage)
}

// RunInSession processes a user message within a named session, replaying history.
// Falls back to stateless Run() if session manager is unavailable.
func (a *Agent) RunInSession(ctx context.Context, sessionKey string, userMessage string) (string, error) {
	originalUserMessage := userMessage

	// Prepend any pending subagent results for this specific session
	if prefix := a.drainPendingResults(sessionKey); prefix != "" {
		userMessage = prefix + "---\nUser message: " + userMessage
	}

	if a.sessions == nil {
		response, err := a.Run(ctx, userMessage)
		if err == nil {
			a.recordMemoryTurn(sessionKey, originalUserMessage, response)
		}
		return response, err
	}

	session, err := a.sessions.Get(sessionKey)
	if err != nil {
		logger.Warn("failed to load session, falling back to stateless", "key", sessionKey, "err", err)
		response, runErr := a.Run(ctx, userMessage)
		if runErr == nil {
			a.recordMemoryTurn(sessionKey, originalUserMessage, response)
		}
		return response, runErr
	}

	// Build messages: system prompt + session history + new user message
	systemPrompt := a.buildSystemPrompt()
	messages := make([]provider.Message, 0, len(session.Messages)+2)
	messages = append(messages, provider.SystemMessage(systemPrompt))
	messages = append(messages, session.Messages...)
	messages = append(messages, provider.UserMessage(userMessage))

	response, err := a.runner.RunWithMessages(ctx, messages)
	if err != nil {
		return "", err
	}

	// Persist the exchange
	session.Messages = append(session.Messages, provider.UserMessage(userMessage))
	session.Messages = append(session.Messages, provider.AssistantMessage(response))

	// Cap history to prevent token overflow.
	const maxSessionMessages = runtimecfg.AgentSessionMaxMessages
	if len(session.Messages) > maxSessionMessages {
		session.Messages = session.Messages[len(session.Messages)-maxSessionMessages:]
	}

	if saveErr := a.sessions.Save(session); saveErr != nil {
		logger.Warn("failed to save session", "key", sessionKey, "err", saveErr)
	}
	a.recordMemoryTurn(sessionKey, originalUserMessage, response)

	return response, nil
}

// buildSystemPrompt builds the system prompt from SOUL.md in workspace.
// SOUL.md should contain the complete system prompt template with placeholders.
func (a *Agent) buildSystemPrompt() string {
	soulPath := filepath.Join(a.workspace, "SOUL.md")
	content, err := os.ReadFile(soulPath)
	if err != nil {
		// Fallback to minimal prompt if SOUL.md doesn't exist
		logger.Warn("SOUL.md not found, using minimal prompt", "path", soulPath)
		return fmt.Sprintf(`You are nagobot, a helpful AI assistant.

Current Time: %s
Workspace: %s
Available Tools: %s
`, time.Now().Format("2006-01-02 15:04 (Monday)"), a.workspace, strings.Join(a.tools.Names(), ", "))
	}

	// Replace placeholders in SOUL.md
	prompt := string(content)
	prompt = strings.ReplaceAll(prompt, "{{TIME}}", time.Now().Format("2006-01-02 15:04 (Monday)"))
	prompt = strings.ReplaceAll(prompt, "{{WORKSPACE}}", a.workspace)
	prompt = strings.ReplaceAll(prompt, "{{TOOLS}}", strings.Join(a.tools.Names(), ", "))

	// User / Agents context files (optional, empty if not present)
	userContent, _ := os.ReadFile(filepath.Join(a.workspace, "USER.md"))
	prompt = strings.ReplaceAll(prompt, "{{USER}}", strings.TrimSpace(string(userContent)))

	agentsContent, _ := os.ReadFile(filepath.Join(a.workspace, "AGENTS.md"))
	prompt = strings.ReplaceAll(prompt, "{{AGENTS}}", strings.TrimSpace(string(agentsContent)))

	// Skills section
	skillsDir := filepath.Join(a.workspace, runtimecfg.WorkspaceSkillsDirName)
	if err := a.skills.ReloadFromDirectory(skillsDir); err != nil {
		logger.Warn("failed to reload skills", "dir", skillsDir, "err", err)
	}
	skillsPrompt := a.skills.BuildPromptSection()
	prompt = strings.ReplaceAll(prompt, "{{SKILLS}}", skillsPrompt)

	// Memory (optional)
	memoryPath := filepath.Join(a.workspace, runtimecfg.WorkspaceMemoryDirName, runtimecfg.MemoryGlobalSummaryFileName)
	memoryContent, _ := os.ReadFile(memoryPath)
	prompt = strings.ReplaceAll(prompt, "{{MEMORY}}", strings.TrimSpace(string(memoryContent)))

	// Today's notes (optional)
	todayFile := time.Now().Format("2006-01-02") + ".md"
	todayPath := filepath.Join(a.workspace, runtimecfg.WorkspaceMemoryDirName, todayFile)
	todayContent, _ := os.ReadFile(todayPath)
	prompt = strings.ReplaceAll(prompt, "{{TODAY}}", strings.TrimSpace(string(todayContent)))

	return prompt
}

func (a *Agent) recordMemoryTurn(sessionKey, userMessage, response string) {
	if a.memory == nil {
		return
	}
	if err := a.memory.RecordTurn(sessionKey, userMessage, response); err != nil {
		logger.Warn("failed to update memory", "session", sessionKey, "err", err)
	}
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
	mu          sync.RWMutex
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
	m.mu.RLock()
	if s, ok := m.cache[key]; ok {
		m.mu.RUnlock()
		return s, nil
	}
	m.mu.RUnlock()

	// Try to load from disk
	path := m.sessionPath(key)
	data, err := os.ReadFile(path)
	if err == nil {
		var s Session
		if err := json.Unmarshal(data, &s); err == nil {
			m.mu.Lock()
			// Another goroutine may have populated the cache while we were reading.
			if cached, ok := m.cache[key]; ok {
				m.mu.Unlock()
				return cached, nil
			}
			m.cache[key] = &s
			m.mu.Unlock()
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
	m.mu.Lock()
	if cached, ok := m.cache[key]; ok {
		m.mu.Unlock()
		return cached, nil
	}
	m.cache[key] = s
	m.mu.Unlock()
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
