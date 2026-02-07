package thread

import (
	"strings"

	"github.com/linanwx/nagobot/provider"
)

// ThreadHook runs during a thread turn with full request context.
type ThreadHook func(ctx HookContext)

// HookContext carries request-building context for hooks.
type HookContext struct {
	ThreadID string
	Type     ThreadType

	SessionKey  string
	SessionPath string

	SessionMessages  []provider.Message
	InjectedMessages []string
	UserMessage      string
	RequestMessages  []provider.Message

	SessionEstimatedTokens int
	RequestEstimatedTokens int
	ContextWindowTokens    int
	ContextWarnRatio       float64
}

// RegisterHook adds a hook for this thread.
func (t *Thread) RegisterHook(h ThreadHook) {
	if h == nil {
		return
	}
	t.mu.Lock()
	t.hooks = append(t.hooks, h)
	t.mu.Unlock()
}

// EnqueueInjectedUserMessage queues a user message to inject before the next real user message.
func (t *Thread) EnqueueInjectedUserMessage(message string) {
	message = strings.TrimSpace(message)
	if message == "" {
		return
	}
	t.mu.Lock()
	t.injectQueue = append(t.injectQueue, message)
	t.mu.Unlock()
}

func (t *Thread) drainInjectQueue() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.injectQueue) == 0 {
		return nil
	}
	out := make([]string, len(t.injectQueue))
	copy(out, t.injectQueue)
	t.injectQueue = nil
	return out
}

func (t *Thread) runHooks(ctx HookContext) {
	t.mu.Lock()
	hooks := make([]ThreadHook, len(t.hooks))
	copy(hooks, t.hooks)
	t.mu.Unlock()

	for _, h := range hooks {
		h(ctx)
	}
}
