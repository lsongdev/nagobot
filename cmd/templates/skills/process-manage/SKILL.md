---
name: process-manage
description: Process management (list, find, kill, monitor resource usage).
tags: [process, system, cross-platform]
---
# Process Management

Monitor and manage system processes. Cross-platform (macOS, Linux).

## List Processes

All processes:
```
exec: ps aux | head -30
```

By name:
```
exec: ps aux | grep -i "PROCESS_NAME" | grep -v grep
```

By PID:
```
exec: ps -p PID -o pid,ppid,user,%cpu,%mem,start,command
```

Tree view (Linux):
```
exec: ps auxf | head -40
```

Tree view (macOS):
```
exec: pstree 2>/dev/null || ps -ef | head -30
```

## Find Process

By name (get PIDs):
```
exec: pgrep -l "PROCESS_NAME"
```

By port:
```
exec: lsof -i :PORT_NUMBER
```

By user:
```
exec: ps -u USERNAME -o pid,%cpu,%mem,command | head -20
```

## Resource Usage

Top CPU consumers:
```
exec: ps aux --sort=-%cpu | head -11
```

Top memory consumers:
```
exec: ps aux --sort=-%mem | head -11
```

System load:
```
exec: uptime
```

One-shot top snapshot:
```
exec: top -bn1 | head -20 2>/dev/null || top -l 1 | head -20
```

## Kill Processes

Kill by PID:
```
exec: kill PID
```

Force kill:
```
exec: kill -9 PID
```

Kill by name:
```
exec: pkill "PROCESS_NAME"
```

Kill all by name:
```
exec: pkill -f "PATTERN"
```

Kill by port:
```
exec: lsof -ti :PORT_NUMBER | xargs kill -9 2>/dev/null && echo "Killed" || echo "No process on port PORT_NUMBER"
```

## Memory & Swap

Linux:
```
exec: free -h
```

macOS:
```
exec: vm_stat | perl -ne '/page size of (\d+)/ and $size=$1; /Pages\s+(\w+):\s+(\d+)/ and printf("%-16s %6.1f MB\n", "$1:", $2 * $size / 1048576)'
```

Cross-platform (Python fallback):
```
exec: python3 -c "import shutil; t,u,f=shutil.disk_usage('/'); print(f'Disk: {u//(1<<30)}GB used / {t//(1<<30)}GB total / {f//(1<<30)}GB free')"
```

## Disk I/O

Linux:
```
exec: iostat -xd 1 2 2>/dev/null | tail -20
```

macOS:
```
exec: iostat -w 1 -c 2 2>/dev/null
```

## Open Files

By process:
```
exec: lsof -p PID | head -30
```

Open file count by process:
```
exec: lsof | awk '{print $1}' | sort | uniq -c | sort -rn | head -15
```

## Service Management (Linux systemd)

Status:
```
exec: systemctl status SERVICE_NAME
```

Start / stop / restart:
```
exec: sudo systemctl start SERVICE_NAME
```
```
exec: sudo systemctl stop SERVICE_NAME
```
```
exec: sudo systemctl restart SERVICE_NAME
```

List running services:
```
exec: systemctl list-units --type=service --state=running
```

## Crontab (System Cron)

List user crontab:
```
exec: crontab -l 2>/dev/null || echo "No crontab for current user"
```

## Notes

- `ps aux` works on both macOS and Linux but output columns may differ slightly.
- `top` flags differ: Linux uses `-b -n 1`, macOS uses `-l 1`.
- `free` is Linux-only; use `vm_stat` or `sysctl` on macOS.
- `systemctl` is Linux-only (systemd); macOS uses `launchctl`.
- `lsof` is available on both platforms and very useful for port/file investigation.
- Some operations (kill, systemctl) may require `sudo`.
