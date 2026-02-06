package thread

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/linanwx/nagobot/provider"
)

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
	m.mu.RLock()
	if s, ok := m.cache[key]; ok {
		m.mu.RUnlock()
		return s, nil
	}
	m.mu.RUnlock()

	path := m.sessionPath(key)
	data, err := os.ReadFile(path)
	if err == nil {
		var s Session
		if err := json.Unmarshal(data, &s); err == nil {
			m.mu.Lock()
			if cached, ok := m.cache[key]; ok {
				m.mu.Unlock()
				return cached, nil
			}
			m.cache[key] = &s
			m.mu.Unlock()
			return &s, nil
		}
	}

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

func (m *SessionManager) sessionPath(key string) string {
	safe := strings.ReplaceAll(key, "/", "_")
	safe = strings.ReplaceAll(safe, ":", "_")
	return filepath.Join(m.sessionsDir, safe+".json")
}
