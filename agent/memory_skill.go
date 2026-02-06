package agent

import (
	"os"
	"path/filepath"

	"github.com/linanwx/nagobot/internal/runtimecfg"
)

const defaultMemorySkillMarkdown = `---
name: memory
description: Retrieve memory via summaries, index, and turn refs.
enabled: true
---
# Memory Skill

Use this skill when you need past context beyond the current turn.

## What Is Already In Context

- Short-term memory: current session history.
- Long-term memory summary: ` + "`memory/MEMORY.md`" + `.
- Today's memory summary: ` + "`memory/YYYY-MM-DD.md`" + `.

## Retrieval Strategy

1. Read summaries first to get the big picture.
2. If details are needed, read the day index file in ` + "`memory/index/*.jsonl`" + `.
3. Match by:
   - ` + "`keywords`" + ` for topic search.
   - ` + "`markers`" + ` (tags like ` + "`#xxx`" + `) for explicit anchors.
   - ` + "`id`" + ` when user references a specific memory ID.
4. Jump to source with ` + "`source_ref`" + ` / ` + "`user_ref`" + ` / ` + "`assistant_ref`" + `.

## Citation Rule

When using retrieved memory, cite the memory ID and ref, e.g.:
` + "`[M-20260206-003] memory/turns/2026-02-06/M-20260206-003.md#L1`" + `.
`

func ensureDefaultMemorySkill(workspace string) error {
	skillsDir := filepath.Join(workspace, runtimecfg.WorkspaceSkillsDirName)
	if err := os.MkdirAll(skillsDir, 0755); err != nil {
		return err
	}

	skillPath := filepath.Join(skillsDir, runtimecfg.MemorySkillFileName)
	if _, err := os.Stat(skillPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	return os.WriteFile(skillPath, []byte(defaultMemorySkillMarkdown), 0644)
}
