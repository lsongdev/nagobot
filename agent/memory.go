package agent

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/linanwx/nagobot/internal/runtimecfg"
)

type memoryStore struct {
	workspace string
	mu        sync.Mutex
}

type memoryIndexEntry struct {
	ID               string   `json:"id"`
	Timestamp        string   `json:"ts"`
	Session          string   `json:"session"`
	SourceFile       string   `json:"source_file"`
	SourceRef        string   `json:"source_ref"`
	UserRef          string   `json:"user_ref"`
	AssistantRef     string   `json:"assistant_ref"`
	UserExcerpt      string   `json:"user_excerpt"`
	AssistantExcerpt string   `json:"assistant_excerpt"`
	Keywords         []string `json:"keywords,omitempty"`
	Markers          []string `json:"markers,omitempty"`
}

var (
	memoryMarkerRe  = regexp.MustCompile(`#([\p{Han}A-Za-z0-9_-]{2,32})`)
	memoryKeywordRe = regexp.MustCompile(`[\p{Han}]{2,}|[A-Za-z][A-Za-z0-9_]{2,}|[0-9]{3,}`)
)

func newMemoryStore(workspace string) *memoryStore {
	return &memoryStore{workspace: workspace}
}

func (m *memoryStore) RecordTurn(sessionKey, userText, assistantText string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	day := now.Format("2006-01-02")
	timestamp := now.Format(time.RFC3339)

	dayEntries, err := m.readDayEntries(day)
	if err != nil {
		return err
	}

	turnID := nextMemoryTurnID(day, dayEntries)
	sourceRel := filepath.ToSlash(filepath.Join(
		runtimecfg.WorkspaceMemoryDirName,
		runtimecfg.MemoryTurnsDirName,
		day,
		turnID+".md",
	))
	sourceAbs := filepath.Join(m.workspace, filepath.FromSlash(sourceRel))

	keywords := extractMemoryKeywords(userText, assistantText)
	markers := extractMemoryMarkers(userText, assistantText)

	entry := memoryIndexEntry{
		ID:               turnID,
		Timestamp:        timestamp,
		Session:          normalizeSessionKey(sessionKey),
		SourceFile:       sourceRel,
		SourceRef:        sourceRel + "#L1",
		UserExcerpt:      firstNChars(userText, runtimecfg.MemoryExcerptChars),
		AssistantExcerpt: firstNChars(assistantText, runtimecfg.MemoryExcerptChars),
		Keywords:         keywords,
		Markers:          markers,
	}

	turnMarkdown, userLine, assistantLine := buildTurnMarkdown(entry, userText, assistantText)
	entry.UserRef = fmt.Sprintf("%s#L%d", sourceRel, userLine)
	entry.AssistantRef = fmt.Sprintf("%s#L%d", sourceRel, assistantLine)

	if err := os.MkdirAll(filepath.Dir(sourceAbs), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(sourceAbs, []byte(turnMarkdown), 0644); err != nil {
		return err
	}

	dayEntries = append(dayEntries, entry)
	dayEntries = m.trimDailyEntries(dayEntries)
	if err := m.writeDayEntries(day, dayEntries); err != nil {
		return err
	}
	if err := m.writeDailySummary(day, dayEntries); err != nil {
		return err
	}
	if err := m.writeGlobalSummary(); err != nil {
		return err
	}

	return nil
}

func (m *memoryStore) trimDailyEntries(dayEntries []memoryIndexEntry) []memoryIndexEntry {
	keep := runtimecfg.MemoryDailyMaxTurns
	if keep <= 0 || len(dayEntries) <= keep {
		return dayEntries
	}

	dropped := dayEntries[:len(dayEntries)-keep]
	for _, d := range dropped {
		if d.SourceFile == "" {
			continue
		}
		_ = os.Remove(filepath.Join(m.workspace, filepath.FromSlash(d.SourceFile)))
	}
	return append([]memoryIndexEntry(nil), dayEntries[len(dayEntries)-keep:]...)
}

func (m *memoryStore) readDayEntries(day string) ([]memoryIndexEntry, error) {
	path := m.dayIndexPath(day)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []memoryIndexEntry{}, nil
		}
		return nil, err
	}
	defer f.Close()

	var entries []memoryIndexEntry
	scanner := bufio.NewScanner(f)
	// Allow long JSON lines.
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var e memoryIndexEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

func (m *memoryStore) writeDayEntries(day string, entries []memoryIndexEntry) error {
	path := m.dayIndexPath(day)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	for _, e := range entries {
		data, err := json.Marshal(e)
		if err != nil {
			return err
		}
		sb.Write(data)
		sb.WriteByte('\n')
	}

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func (m *memoryStore) writeDailySummary(day string, entries []memoryIndexEntry) error {
	summaryPath := filepath.Join(m.workspace, runtimecfg.WorkspaceMemoryDirName, day+".md")
	summary := buildMemorySummary(
		fmt.Sprintf("Daily Memory Summary (%s)", day),
		reverseCopy(entries),
		runtimecfg.MemoryDailySummaryMaxChars,
	)
	return os.WriteFile(summaryPath, []byte(summary), 0644)
}

func (m *memoryStore) writeGlobalSummary() error {
	entries, err := m.collectRecentGlobalEntries(runtimecfg.MemoryGlobalMaxTurns)
	if err != nil {
		return err
	}

	summaryPath := filepath.Join(m.workspace, runtimecfg.WorkspaceMemoryDirName, runtimecfg.MemoryGlobalSummaryFileName)
	summary := buildMemorySummary("Long-term Memory Summary", entries, runtimecfg.MemoryGlobalSummaryMaxChars)
	return os.WriteFile(summaryPath, []byte(summary), 0644)
}

func (m *memoryStore) collectRecentGlobalEntries(limit int) ([]memoryIndexEntry, error) {
	if limit <= 0 {
		return []memoryIndexEntry{}, nil
	}

	indexDir := filepath.Join(m.workspace, runtimecfg.WorkspaceMemoryDirName, runtimecfg.MemoryIndexDirName)
	entries, err := os.ReadDir(indexDir)
	if err != nil {
		if os.IsNotExist(err) {
			return []memoryIndexEntry{}, nil
		}
		return nil, err
	}

	// YYYY-MM-DD sorts lexicographically by time.
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), runtimecfg.MemoryIndexFileExt) {
			continue
		}
		files = append(files, e.Name())
	}
	sortStringsDesc(files)

	var out []memoryIndexEntry
	for _, file := range files {
		day := strings.TrimSuffix(file, runtimecfg.MemoryIndexFileExt)
		dayEntries, err := m.readDayEntries(day)
		if err != nil {
			return nil, err
		}
		for i := len(dayEntries) - 1; i >= 0 && len(out) < limit; i-- {
			out = append(out, dayEntries[i])
		}
		if len(out) >= limit {
			break
		}
	}
	return out, nil
}

func (m *memoryStore) dayIndexPath(day string) string {
	return filepath.Join(
		m.workspace,
		runtimecfg.WorkspaceMemoryDirName,
		runtimecfg.MemoryIndexDirName,
		day+runtimecfg.MemoryIndexFileExt,
	)
}

func buildMemorySummary(title string, entries []memoryIndexEntry, maxChars int) string {
	header := fmt.Sprintf(
		"# %s\n\n- policy: recent turns with excerpt=%d chars\n- indexing: keyword + marker + ref\n\n## Entries (latest first)\n",
		title,
		runtimecfg.MemoryExcerptChars,
	)

	used := runeLen(header)
	var sb strings.Builder
	sb.WriteString(header)
	if len(entries) == 0 {
		sb.WriteString("- (empty)\n")
		return truncateToRunes(sb.String(), maxChars)
	}

	for _, e := range entries {
		line := fmt.Sprintf(
			"- [%s] ref=%s kw=%s mk=%s u=\"%s\" a=\"%s\"\n",
			e.ID,
			e.SourceRef,
			joinOrDash(e.Keywords, "|"),
			joinOrDash(e.Markers, "|"),
			e.UserExcerpt,
			e.AssistantExcerpt,
		)
		lineLen := runeLen(line)
		if used+lineLen > maxChars {
			break
		}
		sb.WriteString(line)
		used += lineLen
	}

	return truncateToRunes(sb.String(), maxChars)
}

func buildTurnMarkdown(entry memoryIndexEntry, userText, assistantText string) (string, int, int) {
	lines := []string{
		"# Turn " + entry.ID,
		"",
		"- timestamp: " + entry.Timestamp,
		"- session: " + entry.Session,
		"- source_ref: " + entry.SourceRef,
		"- keywords: " + joinOrDash(entry.Keywords, ", "),
		"- markers: " + joinOrDash(entry.Markers, ", "),
		"",
	}

	userLine := len(lines) + 1
	lines = append(lines, "## User")
	lines = append(lines, strings.TrimSpace(userText))
	lines = append(lines, "")

	assistantLine := len(lines) + 1
	lines = append(lines, "## Assistant")
	lines = append(lines, strings.TrimSpace(assistantText))
	lines = append(lines, "")

	return strings.Join(lines, "\n"), userLine, assistantLine
}

func nextMemoryTurnID(day string, entries []memoryIndexEntry) string {
	prefix := "M-" + strings.ReplaceAll(day, "-", "") + "-"
	maxSeq := 0
	for _, e := range entries {
		if !strings.HasPrefix(e.ID, prefix) {
			continue
		}
		raw := strings.TrimPrefix(e.ID, prefix)
		n, err := strconv.Atoi(raw)
		if err != nil {
			continue
		}
		if n > maxSeq {
			maxSeq = n
		}
	}
	return fmt.Sprintf("%s%03d", prefix, maxSeq+1)
}

func reverseCopy(in []memoryIndexEntry) []memoryIndexEntry {
	out := make([]memoryIndexEntry, 0, len(in))
	for i := len(in) - 1; i >= 0; i-- {
		out = append(out, in[i])
	}
	return out
}

func firstNChars(s string, n int) string {
	if n <= 0 {
		return ""
	}
	clean := strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
	r := []rune(clean)
	if len(r) <= n {
		return clean
	}
	return string(r[:n])
}

func normalizeSessionKey(sessionKey string) string {
	trimmed := strings.TrimSpace(sessionKey)
	if trimmed == "" {
		return "stateless"
	}
	return trimmed
}

func extractMemoryMarkers(userText, assistantText string) []string {
	all := userText + "\n" + assistantText
	matches := memoryMarkerRe.FindAllStringSubmatch(all, -1)
	seen := make(map[string]struct{})
	var out []string
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		tag := "#" + strings.ToLower(strings.TrimSpace(m[1]))
		if tag == "#" {
			continue
		}
		if _, ok := seen[tag]; ok {
			continue
		}
		seen[tag] = struct{}{}
		out = append(out, tag)
		if len(out) >= runtimecfg.MemoryMaxMarkersPerTurn {
			break
		}
	}
	return out
}

func extractMemoryKeywords(userText, assistantText string) []string {
	all := strings.ToLower(userText + "\n" + assistantText)
	matches := memoryKeywordRe.FindAllString(all, -1)
	seen := make(map[string]struct{})
	var out []string
	for _, token := range matches {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}
		if _, ok := seen[token]; ok {
			continue
		}
		seen[token] = struct{}{}
		out = append(out, token)
		if len(out) >= runtimecfg.MemoryMaxKeywordsPerTurn {
			break
		}
	}
	return out
}

func sortStringsDesc(items []string) {
	for i := 0; i < len(items)-1; i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j] > items[i] {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

func joinOrDash(items []string, sep string) string {
	if len(items) == 0 {
		return "-"
	}
	return strings.Join(items, sep)
}

func runeLen(s string) int {
	return len([]rune(s))
}

func truncateToRunes(s string, max int) string {
	if max <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}
