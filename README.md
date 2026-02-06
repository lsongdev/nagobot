# nagobot

`nagobot` is a lightweight AI assistant built with Go. It supports multiple providers (OpenRouter / Anthropic), tool calling, sessions, cron jobs, multi-channel service mode (CLI / Telegram), and file-based memory indexing.

Inspired by nanobot (`github.com/HKUDS/nanobot`) and openclaw (`github.com/openclaw`).

This project is evolving rapidly.

Repository: <https://github.com/linanwx/nagobot>

## Features

- Providers: `openrouter`, `anthropic`
- Modes:
  - `agent`: single message mode
  - `serve`: long-running service (CLI / Telegram)
- Tools: file read/write, command execution, directory listing, web search/fetch, skill loading, subagents
- Skills:
  - Injects skill summaries by default
  - Loads full skill instructions on demand via `use_skill`
- Memory (file-based):
  - Global summary: `memory/MEMORY.md`
  - Daily summary: `memory/YYYY-MM-DD.md`

## Requirements

- Go `1.23.3+`

## Build

```bash
go build -o nagobot .
```

## Quick Start

1. Initialize config and workspace:

```bash
./nagobot onboard
```

2. Edit config (default: `~/.nagobot/config.json`) and set your API key.

For example:

```json
{
  "providers": {
    "openrouter": {
      "apiKey": "sk-or-v1-xxx"
    }
  }
}
```

3. Start service:

```bash
./nagobot serve
```

## Channels (`serve`)

```bash
# CLI (default)
./nagobot serve

# Enable Telegram
./nagobot serve --telegram

# Enable all configured channels
./nagobot serve --all
```
