---
name: apple-mail
description: Read and manage emails in Apple Mail.
tags: [apple, mail, email, productivity]
---
# Apple Mail

Interact with the macOS Mail app via AppleScript. All commands use `osascript`.

## List All Accounts

```
exec: osascript -e '
tell application "Mail"
    set output to ""
    repeat with acct in accounts
        set output to output & (name of acct) & " (" & (user name of acct) & ")" & linefeed
    end repeat
    return output
end tell
'
```

## List Mailboxes (Folders) for an Account

```
exec: osascript -e '
tell application "Mail"
    set output to ""
    set acct to account "ACCOUNT_NAME"
    repeat with mb in mailboxes of acct
        set msgCount to count of messages in mb
        set output to output & name of mb & " (" & msgCount & " messages)" & linefeed
    end repeat
    return output
end tell
'
```

## Count Unread Messages

```
exec: osascript -e '
tell application "Mail"
    return "Unread: " & (count of (messages of inbox whose read status is false))
end tell
'
```

## List Recent Inbox Messages

List the 10 most recent messages:
```
exec: osascript -e '
tell application "Mail"
    set output to ""
    set msgs to messages 1 thru 10 of inbox
    repeat with msg in msgs
        set msgFrom to sender of msg
        set msgSubject to subject of msg
        set msgDate to date received of msg
        set isRead to read status of msg
        set readMark to "  "
        if not isRead then set readMark to "● "
        set output to output & readMark & (date string of msgDate) & " " & (time string of msgDate) & " | " & msgFrom & " | " & msgSubject & linefeed
    end repeat
    return output
end tell
'
```

## Read a Specific Email by Subject

```
exec: osascript -e '
tell application "Mail"
    set msgs to (messages of inbox whose subject contains "SEARCH_SUBJECT")
    if (count of msgs) > 0 then
        set msg to item 1 of msgs
        set output to "From: " & sender of msg & linefeed
        set output to output & "To: " & (address of to recipient 1 of msg) & linefeed
        set output to output & "Date: " & (date received of msg as text) & linefeed
        set output to output & "Subject: " & subject of msg & linefeed
        set output to output & "---" & linefeed
        set output to output & (content of msg)
        return output
    else
        return "No message found with subject containing: SEARCH_SUBJECT"
    end if
end tell
'
```

## Search Emails

Search by sender:
```
exec: osascript -e '
tell application "Mail"
    set output to ""
    set msgs to (messages of inbox whose sender contains "SENDER_EMAIL")
    repeat with msg in msgs
        set output to output & (date received of msg as text) & " | " & subject of msg & linefeed
        if (count of output) > 5000 then exit repeat
    end repeat
    if output is "" then return "No messages from: SENDER_EMAIL"
    return output
end tell
'
```

## Mark as Read / Unread

```
exec: osascript -e '
tell application "Mail"
    set msgs to (messages of inbox whose subject is "SUBJECT")
    if (count of msgs) > 0 then
        set read status of item 1 of msgs to true
        return "Marked as read."
    end if
    return "Message not found."
end tell
'
```

## Move Message to Mailbox

```
exec: osascript -e '
tell application "Mail"
    set msgs to (messages of inbox whose subject is "SUBJECT")
    if (count of msgs) > 0 then
        set targetMailbox to mailbox "Archive" of account "ACCOUNT_NAME"
        move item 1 of msgs to targetMailbox
        return "Message moved to Archive."
    end if
    return "Message not found."
end tell
'
```

## Create a New Draft

```
exec: osascript -e '
tell application "Mail"
    set newMsg to make new outgoing message with properties {subject:"SUBJECT", content:"EMAIL_BODY", visible:true}
    tell newMsg
        make new to recipient at end of to recipients with properties {address:"recipient@example.com"}
    end tell
    return "Draft created and displayed."
end tell
'
```

## Send an Email

```
exec: osascript -e '
tell application "Mail"
    set newMsg to make new outgoing message with properties {subject:"SUBJECT", content:"EMAIL_BODY", visible:false}
    tell newMsg
        make new to recipient at end of to recipients with properties {address:"recipient@example.com"}
    end tell
    send newMsg
    return "Email sent."
end tell
'
```

## Check for New Mail

```
exec: osascript -e '
tell application "Mail"
    check for new mail
    return "Checking for new mail..."
end tell
'
```

## Notes

- Mail app must be set up with at least one account.
- First run triggers a permission dialog — user must click Allow.
- `ACCOUNT_NAME` should match the name shown in Mail > Settings > Accounts.
- Searching large mailboxes may take a few seconds.
- Sending email requires the account to be properly configured with SMTP.
- Email content may contain HTML. The `content` property returns plain text.
- For security, avoid including sensitive information in exec commands.
