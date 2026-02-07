---
name: compress_context
description: Compress a session JSON file safely when context usage is high, without losing critical state.
---
# Context Compression Skill

Required workflow:
1. Read the session file path provided in the runtime notice.
2. Count file lines first and report the current line count.
3. Generate a compressed session JSON that preserves:
   - top-level fields (`key`, `messages`, `created_at`, `updated_at`)
   - unresolved tasks, decisions, IDs, constraints, pending action items
   - enough recent conversational continuity
4. Write the compressed output to a temporary file path first.
5. Write and run a small validation script that checks:
   - valid JSON
   - required top-level fields exist
   - `messages` is an array
6. Only after validation passes, replace the original session file atomically.

Example compressed session JSON:

```json
{
  "key": "main",
  "messages": [
    {
      "role": "user",
      "content": "[Context Summary]\nThis is a guidance template. The real summary must use the same language as the original conversation and should usually be 1000-3000 characters based on context complexity.\nInclude only high-value context:\n- User preferences and hard constraints\n- Active goals, decisions, and current plan\n- Unresolved issues and pending actions\n- Critical IDs, paths, commands, and references\nRemove repetitive chatter and low-signal turns."
    }
  ],
  "created_at": "2026-02-07T00:00:00Z",
  "updated_at": "2026-02-07T00:10:00Z"
}
```

7. After replacing the file, compression is complete. Continue your original task.
