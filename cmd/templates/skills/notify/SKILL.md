---
name: notify
description: Send macOS native notifications with title, message, and sound.
tags: [macos, notification, utility]
---
# Notify

Send native macOS notifications via AppleScript. Useful for alerting the user when long-running tasks complete, timers fire, or important events occur.

## Send a Basic Notification

```
exec: osascript -e 'display notification "MESSAGE_BODY" with title "TITLE"'
```

## Send with Subtitle

```
exec: osascript -e 'display notification "MESSAGE_BODY" with title "TITLE" subtitle "SUBTITLE"'
```

## Send with Sound

```
exec: osascript -e 'display notification "MESSAGE_BODY" with title "TITLE" sound name "SOUND_NAME"'
```

Available sound names: `Basso`, `Blow`, `Bottle`, `Frog`, `Funk`, `Glass`, `Hero`, `Morse`, `Ping`, `Pop`, `Purr`, `Sosumi`, `Submarine`, `Tink`.

List all system sounds:
```
exec: ls /System/Library/Sounds/
```

## Send Notification and Open URL on Click

Use `terminal-notifier` (if installed via `brew install terminal-notifier`) for richer notifications:
```
exec: terminal-notifier -title "TITLE" -message "MESSAGE" -open "https://example.com"
```

With app icon:
```
exec: terminal-notifier -title "TITLE" -message "MESSAGE" -appIcon /path/to/icon.png
```

## Show a Dialog Box (Blocks Until User Responds)

```
exec: osascript -e 'display dialog "MESSAGE" with title "TITLE" buttons {"OK", "Cancel"} default button "OK"'
```

With text input:
```
exec: osascript -e 'display dialog "Enter value:" default answer "" with title "TITLE"'
```

## Show an Alert

```
exec: osascript -e 'display alert "TITLE" message "DETAIL_MESSAGE" as informational'
```

Alert types: `informational`, `warning`, `critical`.

## Say Text (Text-to-Speech)

```
exec: say "Hello, your task is complete."
```

With specific voice:
```
exec: say -v Samantha "Task complete."
```

List available voices:
```
exec: say -v '?'
```

## Notes

- Notifications appear in macOS Notification Center.
- `display notification` is non-blocking — script returns immediately.
- `display dialog` and `display alert` are blocking — they wait for user response.
- First run may require allowing notifications from the terminal/app in System Settings > Notifications.
- Combine with `manage-cron` skill for scheduled notifications.
