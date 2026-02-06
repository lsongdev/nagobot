---
name: cron
description: Manage scheduled jobs by editing cron.yaml.
tags: [cron, schedule, automation]
enabled: true
---

Use this skill when a user wants to create, update, or inspect timed jobs.

Primary workflow:
1. Read `cron.yaml`.
2. Update the `jobs` list with valid schedule fields.
3. Save changes; the cron service reloads the file automatically.
4. Validate with `go build ./...` and runtime logs.

`cron.yaml` format:

```yaml
jobs:
  - id: "daily-standup"
    name: "Daily Standup Summary"
    enabled: true
    schedule:
      kind: "cron"        # at | every | cron
      expr: "0 9 * * 1-5" # required when kind=cron
      tz: "America/Los_Angeles"
    payload:
      message: "Generate today's standup summary."
      deliver: false
      channel: "telegram" # optional
      to: "123456789"     # optional
    delete_after_run: false

  - id: "run-once"
    name: "One-shot Reminder"
    enabled: true
    schedule:
      kind: "at"
      at_ms: 1760000000000
    payload:
      message: "Run once at exact timestamp."
      deliver: true
      channel: "telegram"
      to: "123456789"
    delete_after_run: true

  - id: "heartbeat"
    name: "Periodic Health Check"
    enabled: true
    schedule:
      kind: "every"
      every_ms: 60000
    payload:
      message: "Check system health every minute."
      deliver: false
```

Schedule rules:
- `kind=at`: requires `at_ms` (unix milliseconds).
- `kind=every`: requires `every_ms > 0`.
- `kind=cron`: requires `expr`; optional `tz`.

Operational notes:
- Prefer stable `id` values for future edits.
- Set `enabled: false` to pause a job without deleting it.
- For one-shot jobs, use `delete_after_run: true` to remove it from in-memory schedules after execution.
