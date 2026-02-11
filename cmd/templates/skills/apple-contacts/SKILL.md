---
name: apple-contacts
description: Search and read contacts from Apple Contacts.
tags: [apple, contacts, productivity]
---
# Apple Contacts

Interact with the macOS Contacts app via AppleScript. All commands use `osascript`.

## List All Contacts

```
exec: osascript -e '
tell application "Contacts"
    set output to ""
    repeat with p in people
        set output to output & (name of p) & linefeed
    end repeat
    return output
end tell
'
```

## List Contact Groups

```
exec: osascript -e '
tell application "Contacts"
    set output to ""
    repeat with g in groups
        set memberCount to count of people in g
        set output to output & name of g & " (" & memberCount & " contacts)" & linefeed
    end repeat
    return output
end tell
'
```

## Search Contact by Name

```
exec: osascript -e '
tell application "Contacts"
    set output to ""
    set results to (every person whose name contains "SEARCH_NAME")
    repeat with p in results
        set output to output & "Name: " & name of p & linefeed
        try
            set output to output & "Phone: " & (value of phone 1 of p) & linefeed
        end try
        try
            set output to output & "Email: " & (value of email 1 of p) & linefeed
        end try
        try
            set output to output & "Company: " & (organization of p) & linefeed
        end try
        set output to output & "---" & linefeed
    end repeat
    if output is "" then return "No contacts found matching: SEARCH_NAME"
    return output
end tell
'
```

## Get Full Contact Details

```
exec: osascript -e '
tell application "Contacts"
    set results to (every person whose name is "FULL_NAME")
    if (count of results) = 0 then return "Contact not found: FULL_NAME"
    set p to item 1 of results

    set output to "Name: " & name of p & linefeed

    try
        set output to output & "Company: " & organization of p & linefeed
    end try
    try
        set output to output & "Title: " & job title of p & linefeed
    end try
    try
        set output to output & "Birthday: " & (birth date of p as text) & linefeed
    end try
    try
        set output to output & "Note: " & note of p & linefeed
    end try

    -- Phones
    repeat with ph in phones of p
        set output to output & "Phone (" & (label of ph) & "): " & (value of ph) & linefeed
    end repeat

    -- Emails
    repeat with em in emails of p
        set output to output & "Email (" & (label of em) & "): " & (value of em) & linefeed
    end repeat

    -- Addresses
    repeat with addr in addresses of p
        set output to output & "Address (" & (label of addr) & "): " & (formatted address of addr) & linefeed
    end repeat

    return output
end tell
'
```

## Search by Phone Number

```
exec: osascript -e '
tell application "Contacts"
    set output to ""
    repeat with p in people
        repeat with ph in phones of p
            if value of ph contains "PHONE_NUMBER" then
                set output to output & name of p & " | " & (value of ph) & linefeed
            end if
        end repeat
    end repeat
    if output is "" then return "No contact found with phone: PHONE_NUMBER"
    return output
end tell
'
```

## Search by Email

```
exec: osascript -e '
tell application "Contacts"
    set output to ""
    repeat with p in people
        repeat with em in emails of p
            if value of em contains "EMAIL_SEARCH" then
                set output to output & name of p & " | " & (value of em) & linefeed
            end if
        end repeat
    end repeat
    if output is "" then return "No contact found with email: EMAIL_SEARCH"
    return output
end tell
'
```

## Create a New Contact

```
exec: osascript -e '
tell application "Contacts"
    set newPerson to make new person with properties {first name:"FIRST_NAME", last name:"LAST_NAME"}
    tell newPerson
        make new phone at end of phones with properties {label:"mobile", value:"PHONE_NUMBER"}
        make new email at end of emails with properties {label:"work", value:"EMAIL_ADDRESS"}
    end tell
    save
    return "Contact created: FIRST_NAME LAST_NAME"
end tell
'
```

## Count Total Contacts

```
exec: osascript -e '
tell application "Contacts"
    return "Total contacts: " & (count of people)
end tell
'
```

## Notes

- Contacts app must be installed (default on macOS).
- First run triggers a permission dialog — user must click Allow.
- Searching by phone is slower on large contact databases (iterates all contacts).
- iCloud-synced contacts are available if signed in.
- Label values: `"mobile"`, `"home"`, `"work"`, `"iPhone"`, etc.
- `save` command is required after creating or modifying contacts.
