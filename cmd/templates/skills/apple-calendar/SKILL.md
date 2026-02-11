---
name: apple-calendar
description: Manage Apple Calendar events (list, create, update, delete).
---
# Apple Calendar

Interact with the macOS Calendar app via AppleScript. All commands use `osascript`.

## List All Calendars

```
exec: osascript -e '
tell application "Calendar"
    set output to ""
    repeat with cal in calendars
        set output to output & (name of cal) & " (" & (uid of cal) & ")" & linefeed
    end repeat
    return output
end tell
'
```

## List Today's Events

```
exec: osascript -e '
set today to current date
set hours of today to 0
set minutes of today to 0
set seconds of today to 0
set tomorrow to today + (1 * days)

tell application "Calendar"
    set output to ""
    repeat with cal in calendars
        set theEvents to (every event of cal whose start date ≥ today and start date < tomorrow)
        repeat with evt in theEvents
            set evtStart to start date of evt
            set evtEnd to end date of evt
            set evtTitle to summary of evt
            set output to output & (time string of evtStart) & " - " & (time string of evtEnd) & " | " & evtTitle & " [" & (name of cal) & "]" & linefeed
        end repeat
    end repeat
    return output
end tell
'
```

## List Upcoming Events (Next N Days)

Replace `7` with desired number of days:
```
exec: osascript -e '
set today to current date
set hours of today to 0
set minutes of today to 0
set seconds of today to 0
set endDate to today + (7 * days)

tell application "Calendar"
    set output to ""
    repeat with cal in calendars
        set theEvents to (every event of cal whose start date ≥ today and start date < endDate)
        repeat with evt in theEvents
            set evtStart to start date of evt
            set evtTitle to summary of evt
            set output to output & (date string of evtStart) & " " & (time string of evtStart) & " | " & evtTitle & " [" & (name of cal) & "]" & linefeed
        end repeat
    end repeat
    return output
end tell
'
```

## Create an Event

Parameters to replace: `CALENDAR_NAME`, `EVENT_TITLE`, `EVENT_LOCATION`, `EVENT_NOTES`, and date/time values.

```
exec: osascript -e '
tell application "Calendar"
    tell calendar "CALENDAR_NAME"
        set startDate to date "2026-02-15 10:00:00"
        set endDate to date "2026-02-15 11:00:00"
        make new event with properties {summary:"EVENT_TITLE", start date:startDate, end date:endDate, location:"EVENT_LOCATION", description:"EVENT_NOTES"}
    end tell
end tell
return "Event created successfully."
'
```

**Important**: Date format depends on the user's system locale. Common formats:
- `"2026/2/15 10:00:00"` (for some Chinese locales)
- `"February 15, 2026 10:00:00 AM"` (English US)
- `"15/02/2026 10:00:00"` (some European locales)

If a specific date format fails, try alternatives or detect locale first:
```
exec: defaults read NSGlobalDomain AppleLocale
```

## Create an All-Day Event

```
exec: osascript -e '
tell application "Calendar"
    tell calendar "CALENDAR_NAME"
        set eventDate to date "2026-02-15 00:00:00"
        make new event with properties {summary:"EVENT_TITLE", start date:eventDate, allday event:true}
    end tell
end tell
return "All-day event created."
'
```

## Delete an Event by Title and Date

```
exec: osascript -e '
set targetDate to date "2026-02-15 00:00:00"
set nextDay to targetDate + (1 * days)

tell application "Calendar"
    repeat with cal in calendars
        set theEvents to (every event of cal whose summary is "EVENT_TITLE" and start date ≥ targetDate and start date < nextDay)
        repeat with evt in theEvents
            delete evt
        end repeat
    end repeat
end tell
return "Event deleted."
'
```

## Search Events by Keyword

```
exec: osascript -e '
set searchTerm to "KEYWORD"
set today to current date
set endDate to today + (30 * days)

tell application "Calendar"
    set output to ""
    repeat with cal in calendars
        set theEvents to (every event of cal whose start date ≥ today and start date < endDate)
        repeat with evt in theEvents
            if summary of evt contains searchTerm then
                set output to output & (date string of (start date of evt)) & " " & (time string of (start date of evt)) & " | " & summary of evt & " [" & (name of cal) & "]" & linefeed
            end if
        end repeat
    end repeat
    return output
end tell
'
```

## Notes

- Calendar app must be installed (default on macOS).
- First run may trigger a permission dialog: "nagobot wants to control Calendar". The user must click Allow.
- Use `tell calendar "CalendarName"` to target a specific calendar. Run "List All Calendars" first to discover names.
- Date parsing is locale-sensitive. Always check the user's locale if date creation fails.
