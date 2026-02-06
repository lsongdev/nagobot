package bus

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/linanwx/nagobot/internal/runtimecfg"
	"github.com/linanwx/nagobot/logger"
)

// SubagentStatus represents the status of a subagent.
type SubagentStatus string

const (
	SubagentStatusPending   SubagentStatus = "pending"
	SubagentStatusRunning   SubagentStatus = "running"
	SubagentStatusCompleted SubagentStatus = "completed"
	SubagentStatusFailed    SubagentStatus = "failed"
)

const completedTaskRetention = 30 * time.Minute

// SubagentOrigin tracks where the spawning request came from,
// so async results can be pushed back to the correct channel/chat.
type SubagentOrigin struct {
	Channel    string // Channel name (e.g., "telegram")
	ReplyTo    string // Chat/user ID for reply routing
	SessionKey string // Session key for context
}

// SubagentTask represents a task for a subagent.
type SubagentTask struct {
	ID          string
	ParentID    string
	Type        string // Agent name (corresponds to agents/*.md file)
	Task        string // The task description/prompt
	Context     string // Additional context
	Origin      SubagentOrigin
	Timeout     time.Duration
	CreatedAt   time.Time
	StartedAt   time.Time
	CompletedAt time.Time
	Status      SubagentStatus
	Result      string
	Error       string
}

// SubagentRunner is a function that runs a subagent task.
type SubagentRunner func(ctx context.Context, task *SubagentTask) (string, error)

// SubagentManager manages subagent lifecycle.
type SubagentManager struct {
	mu            sync.RWMutex
	tasks         map[string]*SubagentTask
	runner        SubagentRunner
	onComplete    func(*SubagentTask) // Called when a task completes or fails
	counter       int64
	maxConcurrent int
	semaphore     chan struct{}
}

// NewSubagentManager creates a new subagent manager.
// onComplete is called (in a goroutine) when a task finishes — may be nil.
func NewSubagentManager(maxConcurrent int, runner SubagentRunner, onComplete func(*SubagentTask)) *SubagentManager {
	if maxConcurrent <= 0 {
		maxConcurrent = runtimecfg.SubagentDefaultMaxConcurrent
	}

	return &SubagentManager{
		tasks:         make(map[string]*SubagentTask),
		runner:        runner,
		onComplete:    onComplete,
		maxConcurrent: maxConcurrent,
		semaphore:     make(chan struct{}, maxConcurrent),
	}
}

// SpawnWithOrigin creates and starts a new subagent task with origin routing info.
func (m *SubagentManager) SpawnWithOrigin(ctx context.Context, parentID, task, taskContext, agentName string, origin SubagentOrigin) (string, error) {
	return m.spawnInternal(ctx, parentID, task, taskContext, agentName, origin)
}

// Spawn creates and starts a new subagent task.
func (m *SubagentManager) Spawn(ctx context.Context, parentID, task, taskContext, agentName string) (string, error) {
	return m.spawnInternal(ctx, parentID, task, taskContext, agentName, SubagentOrigin{})
}

func (m *SubagentManager) spawnInternal(ctx context.Context, parentID, task, taskContext, agentName string, origin SubagentOrigin) (string, error) {
	m.mu.Lock()

	if m.runner == nil {
		m.mu.Unlock()
		return "", fmt.Errorf("subagent runner not configured")
	}

	m.counter++
	idPart := "task"
	if agentName != "" {
		idPart = agentName
	}
	taskID := fmt.Sprintf("sub-%s-%d", idPart, m.counter)

	subTask := &SubagentTask{
		ID:        taskID,
		ParentID:  parentID,
		Type:      agentName,
		Task:      task,
		Context:   taskContext,
		Origin:    origin,
		Timeout:   runtimecfg.SubagentDefaultTimeout,
		CreatedAt: time.Now(),
		Status:    SubagentStatusPending,
	}

	m.tasks[taskID] = subTask
	m.mu.Unlock()

	go m.runTask(ctx, subTask, m.runner)

	logger.Info("subagent spawned", "id", taskID, "agent", agentName, "parent", parentID)
	return taskID, nil
}

// SpawnSync creates and waits for a subagent task to complete.
func (m *SubagentManager) SpawnSync(ctx context.Context, parentID, task, taskContext, agentName string) (string, error) {
	taskID, err := m.Spawn(ctx, parentID, task, taskContext, agentName)
	if err != nil {
		return "", err
	}
	return m.Wait(ctx, taskID)
}

// Wait waits for a task to complete and returns the result.
func (m *SubagentManager) Wait(ctx context.Context, taskID string) (string, error) {
	ticker := time.NewTicker(runtimecfg.SubagentWaitPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-ticker.C:
			m.mu.RLock()
			task, ok := m.tasks[taskID]
			m.mu.RUnlock()

			if !ok {
				return "", fmt.Errorf("task not found: %s", taskID)
			}

			switch task.Status {
			case SubagentStatusCompleted:
				return task.Result, nil
			case SubagentStatusFailed:
				return "", fmt.Errorf("subagent failed: %s", task.Error)
			}
		}
	}
}

// GetTask returns a task by ID.
func (m *SubagentManager) GetTask(taskID string) (*SubagentTask, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	task, ok := m.tasks[taskID]
	return task, ok
}

// runTask executes a subagent task.
func (m *SubagentManager) runTask(ctx context.Context, task *SubagentTask, runner SubagentRunner) {
	// Acquire semaphore
	select {
	case m.semaphore <- struct{}{}:
		defer func() { <-m.semaphore }()
	case <-ctx.Done():
		m.finishTask(task, "", ctx.Err().Error())
		return
	}

	// Update status to running
	m.mu.Lock()
	task.Status = SubagentStatusRunning
	task.StartedAt = time.Now()
	m.mu.Unlock()

	// Create timeout context
	runCtx := ctx
	if task.Timeout > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(ctx, task.Timeout)
		defer cancel()
	}

	// Run the task
	result, err := runner(runCtx, task)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	m.finishTask(task, result, errMsg)
}

// finishTask marks a task as completed or failed and fires the onComplete callback.
func (m *SubagentManager) finishTask(task *SubagentTask, result, errMsg string) {
	m.mu.Lock()
	task.CompletedAt = time.Now()
	if errMsg != "" {
		task.Status = SubagentStatusFailed
		task.Error = errMsg
		logger.Error("subagent failed", "id", task.ID, "err", errMsg)
	} else {
		task.Status = SubagentStatusCompleted
		task.Result = result
		logger.Info("subagent completed", "id", task.ID)
	}
	m.mu.Unlock()

	// Retain completed task status for a while so check_agent can still read it,
	// then release memory.
	taskID := task.ID
	go func() {
		time.Sleep(completedTaskRetention)
		m.mu.Lock()
		delete(m.tasks, taskID)
		m.mu.Unlock()
	}()

	if m.onComplete != nil {
		go m.onComplete(task)
	}
}
