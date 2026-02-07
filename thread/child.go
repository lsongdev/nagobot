package thread

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/linanwx/nagobot/agent"
)

// SpawnChild spawns a child thread for delegated work.
func (t *Thread) SpawnChild(ctx context.Context, ag *agent.Agent, task, taskContext string, wait bool) (string, error) {
	if strings.TrimSpace(task) == "" {
		return "", fmt.Errorf("task is required")
	}

	if ag == nil {
		ag = t.agent
	}
	if ag == nil {
		return "", fmt.Errorf("child agent is not configured")
	}

	childAgent := wrapAgentTaskPlaceholder(ag, task)
	childSessionKey := ""
	if t.cfg != nil && t.cfg.Sessions != nil {
		parentIdentity := strings.TrimSpace(t.sessionKey)
		if parentIdentity == "" {
			parentIdentity = strings.TrimSpace(t.id)
		}
		childSessionKey = generateChildSessionKey(parentIdentity)
	}
	child := newChildWithSession(t.cfg, childAgent, nil, childSessionKey)

	userMessage := task
	if strings.TrimSpace(taskContext) != "" {
		userMessage = fmt.Sprintf("%s\n\nContext:\n%s", task, taskContext)
	}

	if wait {
		return child.Run(ctx, userMessage)
	}

	t.mu.Lock()
	t.childCounter++
	childID := fmt.Sprintf("%s-child-%d", t.id, t.childCounter)
	state := &childState{
		done: make(chan struct{}),
	}
	t.children[childID] = state
	t.mu.Unlock()

	go func() {
		result, err := child.Run(ctx, userMessage)
		t.mu.Lock()
		state.result = result
		state.err = err
		close(state.done)
		t.pendingResults = append(t.pendingResults, pendingChildResult{ID: childID, Result: result, Err: err})
		t.mu.Unlock()
	}()

	return childID, nil
}

// GetChild returns child thread status and result.
func (t *Thread) GetChild(childID string) (status, result string, err error) {
	t.mu.Lock()
	state, ok := t.children[childID]
	if !ok {
		t.mu.Unlock()
		return "", "", fmt.Errorf("child thread not found: %s", childID)
	}

	select {
	case <-state.done:
		result = state.result
		err = state.err
		t.mu.Unlock()
		if err != nil {
			return "failed", "", err
		}
		return "completed", result, nil
	default:
		t.mu.Unlock()
		return "running", "", nil
	}
}

func generateChildSessionKey(parentIdentity string) string {
	parentIdentity = strings.TrimSpace(parentIdentity)
	if parentIdentity == "" {
		parentIdentity = "root"
	}

	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("child:%s:%d", parentIdentity, time.Now().UnixNano())
	}
	return fmt.Sprintf("child:%s:%s", parentIdentity, hex.EncodeToString(buf))
}

func (t *Thread) drainPendingResults() string {
	t.mu.Lock()
	pending := t.pendingResults
	t.pendingResults = nil
	t.mu.Unlock()

	if len(pending) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("[Async thread results since last turn]\n")
	for _, r := range pending {
		if r.Err != nil {
			sb.WriteString(fmt.Sprintf("- %s failed: %v\n", r.ID, r.Err))
			continue
		}
		sb.WriteString(fmt.Sprintf("- %s completed: %s\n", r.ID, r.Result))
	}
	return sb.String()
}

func wrapAgentTaskPlaceholder(base *agent.Agent, task string) *agent.Agent {
	if base == nil {
		return nil
	}
	return &agent.Agent{
		Name:         base.Name,
		ProviderName: base.ProviderName,
		ModelType:    base.ModelType,
		BuildPrompt: func(ctx agent.PromptContext) string {
			if base.BuildPrompt == nil {
				return ""
			}
			prompt := base.BuildPrompt(ctx)
			return strings.ReplaceAll(prompt, "{{TASK}}", task)
		},
	}
}
