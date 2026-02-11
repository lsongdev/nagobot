---
name: apple-notes
description: Manage Apple Notes (list, read, create, search).
---
# Apple Notes

Interact with the macOS Notes app via AppleScript. All commands use `osascript`.

## List All Folders

```
exec: osascript -e '
tell application "Notes"
    set output to ""
    repeat with f in folders
        set noteCount to count of notes in f
        set output to output & name of f & " (" & noteCount & " notes)" & linefeed
    end repeat
    return output
end tell
'
```

## List Notes in a Folder

Replace `FOLDER_NAME` with the target folder (use `"Notes"` for the default folder):
```
exec: osascript -e '
tell application "Notes"
    set output to ""
    set theNotes to notes of folder "FOLDER_NAME"
    repeat with n in theNotes
        set output to output & name of n & " | " & (modification date of n as text) & linefeed
    end repeat
    return output
end tell
'
```

## List All Notes (All Folders)

```
exec: osascript -e '
tell application "Notes"
    set output to ""
    repeat with f in folders
        set folderName to name of f
        repeat with n in notes of f
            set output to output & "[" & folderName & "] " & name of n & " | " & (modification date of n as text) & linefeed
        end repeat
    end repeat
    return output
end tell
'
```

## Read a Note by Name

```
exec: osascript -e '
tell application "Notes"
    set theNotes to notes whose name is "NOTE_TITLE"
    if (count of theNotes) > 0 then
        set theNote to item 1 of theNotes
        set noteName to name of theNote
        set noteBody to plaintext of theNote
        return "Title: " & noteName & linefeed & linefeed & noteBody
    else
        return "Note not found: NOTE_TITLE"
    end if
end tell
'
```

## Create a New Note

Replace `FOLDER_NAME`, `NOTE_TITLE`, and `NOTE_BODY`. The body supports HTML for rich text.

Plain text note:
```
exec: osascript -e '
tell application "Notes"
    tell folder "FOLDER_NAME"
        make new note with properties {name:"NOTE_TITLE", body:"NOTE_BODY"}
    end tell
end tell
return "Note created successfully."
'
```

Rich text note (HTML body):
```
exec: osascript -e '
tell application "Notes"
    tell folder "Notes"
        make new note with properties {name:"NOTE_TITLE", body:"<h1>Title</h1><p>Paragraph text here.</p><ul><li>Item 1</li><li>Item 2</li></ul>"}
    end tell
end tell
return "Note created successfully."
'
```

## Append to an Existing Note

```
exec: osascript -e '
tell application "Notes"
    set theNotes to notes whose name is "NOTE_TITLE"
    if (count of theNotes) > 0 then
        set theNote to item 1 of theNotes
        set currentBody to body of theNote
        set body of theNote to currentBody & "<br><p>APPENDED_TEXT</p>"
        return "Text appended successfully."
    else
        return "Note not found: NOTE_TITLE"
    end if
end tell
'
```

## Search Notes by Keyword

```
exec: osascript -e '
set searchTerm to "KEYWORD"
tell application "Notes"
    set output to ""
    repeat with f in folders
        repeat with n in notes of f
            if (name of n contains searchTerm) or (plaintext of n contains searchTerm) then
                set output to output & "[" & (name of f) & "] " & name of n & linefeed
            end if
        end repeat
    end repeat
    if output is "" then
        return "No notes found matching: " & searchTerm
    end if
    return output
end tell
'
```

## Delete a Note by Name

```
exec: osascript -e '
tell application "Notes"
    set theNotes to notes whose name is "NOTE_TITLE"
    if (count of theNotes) > 0 then
        delete item 1 of theNotes
        return "Note deleted."
    else
        return "Note not found: NOTE_TITLE"
    end if
end tell
'
```

## Notes

- Notes app must be installed (default on macOS).
- First run may trigger a permission dialog. The user must click Allow.
- Default folder is typically named `"Notes"` (or localized equivalent, e.g., `"备忘录"` in Chinese).
- Note body uses HTML internally. Use `plaintext` property to get clean text.
- Very large notes may be truncated by exec output limits.
- iCloud-synced notes are accessible if the user is signed in.
