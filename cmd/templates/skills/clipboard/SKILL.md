---
name: clipboard
description: Read from and write to the macOS clipboard.
tags: [macos, clipboard, utility]
---
# Clipboard

Read and write the macOS system clipboard using `pbcopy` and `pbpaste`.

## Read Clipboard Contents

```
exec: pbpaste
```

## Write Text to Clipboard

```
exec: echo -n "TEXT_TO_COPY" | pbcopy
```

Write file contents to clipboard:
```
exec: cat /path/to/file.txt | pbcopy
```

## Write Command Output to Clipboard

```
exec: date | pbcopy
```

```
exec: pwd | pbcopy
```

## Pipe Between Clipboard and Tools

Clipboard → process → clipboard:
```
exec: pbpaste | sort | pbcopy
```

Clipboard → file:
```
exec: pbpaste > /path/to/output.txt
```

File → clipboard:
```
exec: pbcopy < /path/to/input.txt
```

## Copy Image to Clipboard (via AppleScript)

```
exec: osascript -e '
set theFile to POSIX file "/path/to/image.png"
set the clipboard to (read theFile as «class PNGf»)
return "Image copied to clipboard."
'
```

## Get Clipboard Type Info

```
exec: osascript -e '
clipboard info
'
```

## Clear Clipboard

```
exec: osascript -e 'set the clipboard to ""'
```

## Notes

- `pbcopy` reads from stdin and places content on the clipboard.
- `pbpaste` writes clipboard content to stdout.
- These are macOS-only commands (not available on Linux).
- Binary clipboard contents (images, etc.) require AppleScript.
- Clipboard changes are immediate and system-wide.
