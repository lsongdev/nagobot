---
name: apple-reminders
description: Manage Apple Reminders (list, create, complete, delete).
---
# Apple Reminders

Interact with the macOS Reminders app via AppleScript. All commands use `osascript`.

## List All Reminder Lists

```
exec: osascript -e '
tell application "Reminders"
    set output to ""
    repeat with lst in lists
        set reminderCount to count of reminders in lst
        set incompleteCount to count of (reminders of lst whose completed is false)
        set output to output & name of lst & " (" & incompleteCount & " incomplete / " & reminderCount & " total)" & linefeed
    end repeat
    return output
end tell
'
```

## List Incomplete Reminders in a List

Replace `LIST_NAME` with the target list (e.g., `"Reminders"` or `"提醒事项"`):
```
exec: osascript -e '
tell application "Reminders"
    set output to ""
    set theReminders to (reminders of list "LIST_NAME" whose completed is false)
    repeat with r in theReminders
        set rName to name of r
        set rDate to ""
        try
            set rDate to " | Due: " & (due date of r as text)
        end try
        set rPriority to priority of r
        set priLabel to ""
        if rPriority is 1 then set priLabel to " ⚡"
        if rPriority is 5 then set priLabel to " ❗"
        if rPriority is 9 then set priLabel to " ↓"
        set output to output & "- " & rName & rDate & priLabel & linefeed
    end repeat
    if output is "" then
        return "No incomplete reminders in list: LIST_NAME"
    end if
    return output
end tell
'
```

## List All Incomplete Reminders (All Lists)

```
exec: osascript -e '
tell application "Reminders"
    set output to ""
    repeat with lst in lists
        set theReminders to (reminders of lst whose completed is false)
        if (count of theReminders) > 0 then
            set output to output & "=== " & name of lst & " ===" & linefeed
            repeat with r in theReminders
                set rName to name of r
                set rDate to ""
                try
                    set rDate to " | Due: " & (due date of r as text)
                end try
                set output to output & "- " & rName & rDate & linefeed
            end repeat
            set output to output & linefeed
        end if
    end repeat
    if output is "" then
        return "No incomplete reminders found."
    end if
    return output
end tell
'
```

## Create a Reminder

Simple reminder (no due date):
```
exec: osascript -e '
tell application "Reminders"
    tell list "LIST_NAME"
        make new reminder with properties {name:"REMINDER_TITLE", body:"REMINDER_NOTES"}
    end tell
end tell
return "Reminder created."
'
```

Reminder with due date:
```
exec: osascript -e '
tell application "Reminders"
    tell list "LIST_NAME"
        set dueDate to date "2026-02-15 09:00:00"
        make new reminder with properties {name:"REMINDER_TITLE", body:"REMINDER_NOTES", due date:dueDate}
    end tell
end tell
return "Reminder created with due date."
'
```

Reminder with priority (1=high, 5=medium, 9=low):
```
exec: osascript -e '
tell application "Reminders"
    tell list "LIST_NAME"
        make new reminder with properties {name:"REMINDER_TITLE", priority:1}
    end tell
end tell
return "High-priority reminder created."
'
```

## Complete a Reminder

```
exec: osascript -e '
tell application "Reminders"
    set theReminders to (reminders whose name is "REMINDER_TITLE" and completed is false)
    if (count of theReminders) > 0 then
        set completed of item 1 of theReminders to true
        return "Reminder marked as completed."
    else
        return "Reminder not found or already completed: REMINDER_TITLE"
    end if
end tell
'
```

## Delete a Reminder

```
exec: osascript -e '
tell application "Reminders"
    set theReminders to (reminders whose name is "REMINDER_TITLE" and completed is false)
    if (count of theReminders) > 0 then
        delete item 1 of theReminders
        return "Reminder deleted."
    else
        return "Reminder not found: REMINDER_TITLE"
    end if
end tell
'
```

## Search Reminders by Keyword

```
exec: osascript -e '
set searchTerm to "KEYWORD"
tell application "Reminders"
    set output to ""
    repeat with lst in lists
        set theReminders to (reminders of lst whose completed is false)
        repeat with r in theReminders
            if name of r contains searchTerm then
                set rDate to ""
                try
                    set rDate to " | Due: " & (due date of r as text)
                end try
                set output to output & "[" & (name of lst) & "] " & name of r & rDate & linefeed
            end if
        end repeat
    end repeat
    if output is "" then
        return "No reminders found matching: " & searchTerm
    end if
    return output
end tell
'
```

## Create a New Reminder List

```
exec: osascript -e '
tell application "Reminders"
    make new list with properties {name:"NEW_LIST_NAME"}
end tell
return "Reminder list created."
'
```

## Notes

- Reminders app must be installed (default on macOS).
- First run may trigger a permission dialog. The user must click Allow.
- Default list is typically `"Reminders"` (or localized: `"提醒事项"` in Chinese).
- Date format is locale-sensitive. Check locale with `defaults read NSGlobalDomain AppleLocale` if date creation fails.
- Priority values: 0 = none, 1 = high, 5 = medium, 9 = low.
- iCloud-synced reminders are accessible if signed in.
- Completed reminders can be listed by changing `completed is false` to `completed is true`.
