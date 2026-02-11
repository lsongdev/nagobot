// Package msg defines the WakeMessage type shared between thread and tools.
package msg

import "context"

// Sink defines how thread output is delivered.
type Sink func(ctx context.Context, response string) error

// ThreadInfo holds the summary status of a thread.
type ThreadInfo struct {
	ID         string `json:"id"`
	SessionKey string `json:"sessionKey"`
	State      string `json:"state"`   // "running", "pending", "idle"
	Pending    int    `json:"pending"`
}

// WakeMessage is an item in a thread's wake queue.
type WakeMessage struct {
	Source    string            // Wake source: "telegram", "cron", "child_completed", etc.
	Message   string            // Wake payload text.
	Sink      Sink              // Per-wake response delivery (nil = drop response).
	SinkLabel string            // Descriptive label for this sink.
	AgentName string            // Optional agent name override for this wake.
	Vars      map[string]string // Optional vars override for this wake.
}
