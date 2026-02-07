#!/usr/bin/env bash
set -euo pipefail

CONFIG_DIR="${1:-$HOME/.nagobot}"
CONFIG_PATH="$CONFIG_DIR/config.yaml"
BACKUP_PATH=""
RESTORED=0

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "$SCRIPT_DIR/.." && pwd)"

restore_config() {
  if [[ -n "$BACKUP_PATH" && -f "$BACKUP_PATH" && "$RESTORED" -eq 0 ]]; then
    mkdir -p "$CONFIG_DIR"
    cp "$BACKUP_PATH" "$CONFIG_PATH"
    RESTORED=1
    echo "Restored config to: $CONFIG_PATH"
  fi
  if [[ -n "$BACKUP_PATH" && -f "$BACKUP_PATH" ]]; then
    rm -f "$BACKUP_PATH"
  fi
}

trap restore_config EXIT

if [[ ! -f "$CONFIG_PATH" ]]; then
  echo "Config not found: $CONFIG_PATH" >&2
  exit 1
fi

if ! command -v trash >/dev/null 2>&1; then
  echo "'trash' command not found. Install it first (e.g. brew install trash)." >&2
  exit 1
fi

BACKUP_PATH="$(mktemp "${TMPDIR:-/tmp}/nagobot-config-backup-XXXXXX.yaml")"
cp "$CONFIG_PATH" "$BACKUP_PATH"
echo "Backed up config to: $BACKUP_PATH"

echo "Trashing config directory: $CONFIG_DIR"
trash "$CONFIG_DIR"

echo "Running onboard in repo: $REPO_ROOT"
cd "$REPO_ROOT"
go run . onboard

# Restore immediately after onboard; trap also protects failure paths.
restore_config
trap - EXIT

echo "Done."
