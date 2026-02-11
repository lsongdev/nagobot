---
name: screenshot
description: Capture screenshots on macOS (full screen, window, or region).
tags: [macos, screenshot, utility]
---
# Screenshot

Capture screenshots using macOS built-in `screencapture` command.

## Capture Full Screen to File

```
exec: screencapture -x {{WORKSPACE}}/.tmp/screenshot.png
```

The `-x` flag suppresses the camera shutter sound.

## Capture Full Screen to Clipboard

```
exec: screencapture -xc
```

## Capture Selected Region (Interactive)

Prompts the user to select a screen region:
```
exec: screencapture -xi {{WORKSPACE}}/.tmp/region.png
```

## Capture a Specific Window (Interactive)

Prompts the user to click a window:
```
exec: screencapture -xw {{WORKSPACE}}/.tmp/window.png
```

## Capture with Delay (Seconds)

Useful to set up the screen before capture:
```
exec: screencapture -x -T 5 {{WORKSPACE}}/.tmp/delayed.png
```

## Capture Specific Display

List displays and capture by display number:
```
exec: screencapture -x -D 1 {{WORKSPACE}}/.tmp/display1.png
```

## Capture as JPEG (with Quality)

```
exec: screencapture -x -t jpg {{WORKSPACE}}/.tmp/screenshot.jpg
```

## Capture as PDF

```
exec: screencapture -x -t pdf {{WORKSPACE}}/.tmp/screenshot.pdf
```

## Capture and Open Immediately in Preview

```
exec: screencapture -x -P {{WORKSPACE}}/.tmp/screenshot.png
```

## Open Screenshot in Default Viewer

After capturing:
```
exec: open {{WORKSPACE}}/.tmp/screenshot.png
```

## Flag Reference

- `-x` — no sound
- `-c` — capture to clipboard
- `-w` — capture window (interactive click)
- `-i` — capture interactive selection
- `-T <seconds>` — delay before capture
- `-D <display>` — capture specific display
- `-t <format>` — format: `png` (default), `jpg`, `pdf`, `tiff`
- `-P` — open in Preview after capture
- `-R <x,y,w,h>` — capture specific rectangle

## Notes

- Screenshots are saved to `{{WORKSPACE}}/.tmp/` by default. Change the path as needed.
- Non-interactive captures (`-x` without `-i` or `-w`) work in headless/background contexts.
- First use may require Screen Recording permission in System Settings > Privacy & Security > Screen Recording.
- For annotating screenshots, use `open -a Preview <file>` after capture.
