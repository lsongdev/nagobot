package agent

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/linanwx/nagobot/logger"
)

// AgentDef represents an agent template file under workspace/agents.
type AgentDef struct {
	Name string // Filename without .md extension
	Path string // Full path to the template file
}

// AgentRegistry loads agent templates from workspace/agents.
type AgentRegistry struct {
	workspace string
	agentsDir string
	agents    map[string]*AgentDef
	mu        sync.RWMutex
}

// NewRegistry creates a registry from workspace/agents and loads all templates.
func NewRegistry(workspace string) *AgentRegistry {
	r := &AgentRegistry{
		workspace: workspace,
		agentsDir: filepath.Join(workspace, "agents"),
		agents:    make(map[string]*AgentDef),
	}
	r.load()
	return r
}

func (r *AgentRegistry) load() {
	entries, err := os.ReadDir(r.agentsDir)
	if err != nil {
		logger.Debug("agents directory not found", "dir", r.agentsDir)
		return
	}

	next := make(map[string]*AgentDef)
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		next[name] = &AgentDef{
			Name: name,
			Path: filepath.Join(r.agentsDir, entry.Name()),
		}
	}

	r.mu.Lock()
	r.agents = next
	r.mu.Unlock()
}

// Get returns a prompt-builder agent by name, or nil if not found.
func (r *AgentRegistry) Get(name string) *Agent {
	r.mu.RLock()
	def, ok := r.agents[name]
	r.mu.RUnlock()
	if !ok {
		return nil
	}

	return NewTemplateAgent(def.Name, def.Path, r.workspace)
}

// List returns all available agent names.
func (r *AgentRegistry) List() []string {
	r.mu.RLock()
	names := make([]string, 0, len(r.agents))
	for name := range r.agents {
		names = append(names, name)
	}
	r.mu.RUnlock()
	sort.Strings(names)
	return names
}

// Reload reloads agent templates from disk.
func (r *AgentRegistry) Reload() {
	r.load()
}
