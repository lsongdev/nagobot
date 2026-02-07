package thread

import (
	"encoding/json"
	"errors"
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

	s, err := m.loadFromDisk(key)
	if err != nil {
		return nil, err
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

// Reload forces loading session state from disk and refreshes cache.
func (m *SessionManager) Reload(key string) (*Session, error) {
	s, err := m.loadFromDisk(key)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
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

// PathForKey returns the on-disk session file path for a session key.
func (m *SessionManager) PathForKey(key string) string {
	return m.sessionPath(key)
}

func (m *SessionManager) loadFromDisk(key string) (*Session, error) {
	path := m.sessionPath(key)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			now := time.Now()
			return &Session{
				Key:       key,
				Messages:  []provider.Message{},
				CreatedAt: now,
				UpdatedAt: now,
			}, nil
		}
		return nil, err
	}

	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	if strings.TrimSpace(s.Key) == "" {
		s.Key = key
	}
	if s.Messages == nil {
		s.Messages = []provider.Message{}
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	if s.UpdatedAt.IsZero() {
		s.UpdatedAt = s.CreatedAt
	}
	return &s, nil
}
