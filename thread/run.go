package thread

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/linanwx/nagobot/agent"
	"github.com/linanwx/nagobot/internal/runtimecfg"
	"github.com/linanwx/nagobot/logger"
	"github.com/linanwx/nagobot/provider"
	"github.com/linanwx/nagobot/tools"
)

// Run executes one thread turn.
func (t *Thread) Run(ctx context.Context, userMessage string) (string, error) {
	prov, err := t.resolveProvider()
	if err != nil {
		return "", err
	}

	runtimeTools := t.runtimeTools()
	skillsSection := t.buildSkillsSection()

	if prefix := t.drainPendingResults(); prefix != "" {
		userMessage = prefix + "---\nUser message: " + userMessage
	}

	promptCtx := agent.PromptContext{
		Workspace: t.workspace,
		Time:      time.Now(),
		ToolNames: runtimeTools.Names(),
		Skills:    skillsSection,
	}

	systemPrompt := ""
	if t.agent != nil && t.agent.BuildPrompt != nil {
		systemPrompt = t.agent.BuildPrompt(promptCtx)
	}
	if strings.TrimSpace(systemPrompt) == "" {
		systemPrompt = agent.NewRawAgent("fallback", "You are a helpful AI assistant.").BuildPrompt(promptCtx)
	}

	messages := make([]provider.Message, 0, 2)
	messages = append(messages, provider.SystemMessage(systemPrompt))

	session := t.loadSession()
	if session != nil {
		messages = append(messages, session.Messages...)
	}

	messages = append(messages, provider.UserMessage(userMessage))

	sessionEstimatedTokens := 0
	if session != nil {
		sessionEstimatedTokens = estimateMessagesTokens(session.Messages)
	}
	requestEstimatedTokens := estimateMessagesTokens(messages)
	contextWindowTokens, contextWarnRatio := t.contextBudget()
	logger.Debug(
		"context estimate",
		"threadID", t.id,
		"threadType", t.kind,
		"sessionKey", t.sessionKey,
		"sessionEstimatedTokens", sessionEstimatedTokens,
		"requestEstimatedTokens", requestEstimatedTokens,
		"contextWindowTokens", contextWindowTokens,
		"contextWarnRatio", contextWarnRatio,
	)

	if sessionPath, ok := t.sessionFilePath(); ok {
		threshold := int(float64(contextWindowTokens) * contextWarnRatio)
		if threshold <= 0 {
			threshold = contextWindowTokens
		}
		if requestEstimatedTokens >= threshold {
			usageRatio := float64(requestEstimatedTokens) / float64(contextWindowTokens)
			notice := t.buildCompressionNotice(requestEstimatedTokens, contextWindowTokens, usageRatio, sessionPath)
			userMessage = notice + "\n\n---\nUser message: " + userMessage
			messages[len(messages)-1] = provider.UserMessage(userMessage)
			requestEstimatedTokens = estimateMessagesTokens(messages)

			logger.Info(
				"context threshold reached, compression reminder injected",
				"threadID", t.id,
				"threadType", t.kind,
				"sessionKey", t.sessionKey,
				"sessionPath", sessionPath,
				"requestEstimatedTokens", requestEstimatedTokens,
				"contextWindowTokens", contextWindowTokens,
				"thresholdTokens", threshold,
			)
		}
	}

	runner := NewRunner(prov, runtimeTools)
	response, err := runner.RunWithMessages(ctx, messages)
	if err != nil {
		return "", err
	}

	if session != nil {
		session.Messages = append(session.Messages, provider.UserMessage(userMessage))
		session.Messages = append(session.Messages, provider.AssistantMessage(response))

		if saveErr := t.cfg.Sessions.Save(session); saveErr != nil {
			logger.Warn("failed to save session", "key", t.sessionKey, "err", saveErr)
		}
	}

	if t.sink != nil {
		if err := t.sink(ctx, response); err != nil {
			return "", err
		}
	}

	return response, nil
}

func (t *Thread) resolveProvider() (provider.Provider, error) {
	if t.agent != nil && (strings.TrimSpace(t.agent.ProviderName) != "" || strings.TrimSpace(t.agent.ModelType) != "") {
		if t.cfg.ProviderFactory == nil {
			return nil, fmt.Errorf("provider override requested but provider factory is not configured")
		}
		return t.cfg.ProviderFactory(t.agent.ProviderName, t.agent.ModelType)
	}

	if t.provider == nil {
		if t.cfg.ProviderFactory == nil {
			return nil, fmt.Errorf("default provider is not configured")
		}
		return t.cfg.ProviderFactory("", "")
	}

	return t.provider, nil
}

func (t *Thread) runtimeTools() *tools.Registry {
	runtimeTools := tools.NewRegistry()
	if t.tools != nil {
		runtimeTools = t.tools.Clone()
	}

	if t.allowSpawn {
		runtimeTools.Register(tools.NewSpawnThreadTool(t, t.agents))
		runtimeTools.Register(tools.NewCheckThreadTool(t))
	}

	return runtimeTools
}

func (t *Thread) loadSession() *Session {
	if t.cfg == nil || strings.TrimSpace(t.sessionKey) == "" || t.cfg.Sessions == nil {
		return nil
	}

	session, err := t.cfg.Sessions.Get(t.sessionKey)
	if err != nil {
		logger.Warn("failed to load session", "key", t.sessionKey, "err", err)
		return nil
	}
	return session
}

func (t *Thread) buildSkillsSection() string {
	if t.skills == nil || strings.TrimSpace(t.workspace) == "" {
		return ""
	}

	skillsDir := filepath.Join(t.workspace, runtimecfg.WorkspaceSkillsDirName)
	if err := t.skills.ReloadFromDirectory(skillsDir); err != nil {
		logger.Warn("failed to reload skills", "dir", skillsDir, "err", err)
	}
	return t.skills.BuildPromptSection()
}
