---
name: http-request
description: Make HTTP requests with curl (GET, POST, headers, auth, file upload).
tags: [http, api, curl, cross-platform]
---
# HTTP Request

Make HTTP requests using `curl`. Cross-platform (macOS, Linux, Windows WSL).

## GET Request

Basic:
```
exec: curl -s "https://api.example.com/endpoint"
```

With headers:
```
exec: curl -s -H "Accept: application/json" -H "Authorization: Bearer TOKEN" "https://api.example.com/endpoint"
```

Show response headers:
```
exec: curl -sI "https://example.com"
```

Show both headers and body:
```
exec: curl -si "https://example.com"
```

## POST Request

JSON body:
```
exec: curl -s -X POST "https://api.example.com/endpoint" -H "Content-Type: application/json" -d '{"key":"value","name":"test"}'
```

Form data:
```
exec: curl -s -X POST "https://api.example.com/endpoint" -d "field1=value1&field2=value2"
```

## PUT / PATCH / DELETE

```
exec: curl -s -X PUT "https://api.example.com/resource/1" -H "Content-Type: application/json" -d '{"key":"updated"}'
```

```
exec: curl -s -X PATCH "https://api.example.com/resource/1" -H "Content-Type: application/json" -d '{"field":"patched"}'
```

```
exec: curl -s -X DELETE "https://api.example.com/resource/1"
```

## Authentication

Bearer token:
```
exec: curl -s -H "Authorization: Bearer TOKEN" "https://api.example.com/me"
```

Basic auth:
```
exec: curl -s -u "username:password" "https://api.example.com/endpoint"
```

API key in header:
```
exec: curl -s -H "X-API-Key: KEY" "https://api.example.com/endpoint"
```

## File Upload

Multipart:
```
exec: curl -s -X POST "https://api.example.com/upload" -F "file=@/path/to/file.pdf" -F "description=My file"
```

## File Download

Save to file:
```
exec: curl -sL -o /path/to/output.zip "https://example.com/file.zip"
```

With progress bar:
```
exec: curl -L -o /path/to/output.zip "https://example.com/file.zip"
```

## Parse JSON Response (with jq)

Extract field:
```
exec: curl -s "https://api.example.com/data" | jq '.results[0].name'
```

Pretty print:
```
exec: curl -s "https://api.example.com/data" | jq .
```

Filter array:
```
exec: curl -s "https://api.example.com/items" | jq '[.[] | select(.status == "active")]'
```

## Timing & Debug

Show timing info:
```
exec: curl -s -o /dev/null -w "HTTP %{http_code} | Time: %{time_total}s | Size: %{size_download} bytes\n" "https://example.com"
```

Verbose (debug):
```
exec: curl -v "https://example.com" 2>&1 | head -30
```

## Follow Redirects

```
exec: curl -sL "https://short.url/abc"
```

## Common API Examples

GitHub API:
```
exec: curl -s -H "Authorization: token GITHUB_TOKEN" "https://api.github.com/user/repos?per_page=5" | jq '.[].full_name'
```

Check IP:
```
exec: curl -s ifconfig.me
```

## Notes

- `-s` = silent (no progress bar), `-S` = show errors even in silent mode.
- `-L` = follow redirects.
- `-o FILE` = save output to file.
- `-w FORMAT` = write-out format for metrics.
- `jq` is recommended for JSON parsing; install via package manager if not available.
- For APIs requiring complex auth flows (OAuth), consider scripting the token exchange.
- Works identically on macOS, Linux, and Windows (WSL/Git Bash).
