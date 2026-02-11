---
name: system-info
description: Query macOS system info (disk, memory, CPU, network, battery).
tags: [macos, system, diagnostics, utility]
---
# System Info

Query macOS system information using built-in commands.

## Disk Usage

Overview of all volumes:
```
exec: df -h
```

Directory size breakdown:
```
exec: du -sh ~/* 2>/dev/null | sort -hr | head -20
```

Specific directory:
```
exec: du -sh /path/to/directory
```

## Memory Usage

```
exec: vm_stat | perl -ne '/page size of (\d+)/ and $size=$1; /Pages\s+(\w+):\s+(\d+)/ and printf("%-16s %6.1f MB\n", "$1:", $2 * $size / 1048576)'
```

Quick memory summary:
```
exec: sysctl hw.memsize | awk '{printf "Total RAM: %.1f GB\n", $2/1073741824}'
```

Top memory consumers:
```
exec: ps aux --sort=-%mem | head -11
```

## CPU Info

CPU model and core count:
```
exec: sysctl -n machdep.cpu.brand_string && echo "Cores: $(sysctl -n hw.ncpu) ($(sysctl -n hw.physicalcpu) physical)"
```

Current CPU load:
```
exec: top -l 1 -n 0 | grep "CPU usage"
```

Top CPU consumers:
```
exec: ps aux --sort=-%cpu | head -11
```

## System Overview

macOS version and hardware:
```
exec: sw_vers && echo "---" && system_profiler SPHardwareDataType 2>/dev/null | grep -E "Model|Chip|Memory|Serial"
```

Uptime:
```
exec: uptime
```

Hostname and user:
```
exec: echo "Hostname: $(hostname)" && echo "User: $(whoami)"
```

## Network Info

Active network interfaces:
```
exec: ifconfig | grep -E "^[a-z]|inet " | grep -B1 "inet "
```

Public IP:
```
exec: curl -s ifconfig.me
```

Wi-Fi info:
```
exec: /System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I | grep -E "SSID|BSSID|channel|RSSI"
```

Network speed test (simple):
```
exec: curl -s -o /dev/null -w "Download speed: %{speed_download} bytes/sec\nTime: %{time_total}s\n" https://speed.cloudflare.com/__down?bytes=10000000
```

DNS servers:
```
exec: scutil --dns | grep nameserver | head -5
```

Active connections:
```
exec: netstat -an | grep ESTABLISHED | wc -l
```

## Battery Status (Laptops)

```
exec: pmset -g batt
```

Detailed battery info:
```
exec: system_profiler SPPowerDataType 2>/dev/null | grep -E "Charge|Capacity|Cycle|Condition"
```

## Storage Devices

```
exec: diskutil list
```

Volume info:
```
exec: diskutil info /
```

## Running Processes

Total processes:
```
exec: ps aux | wc -l
```

Find process by name:
```
exec: pgrep -l "PROCESS_NAME"
```

Kill process by name:
```
exec: pkill "PROCESS_NAME"
```

## Thermal & Fan (Intel Macs)

```
exec: sudo powermetrics --samplers smc -i 1 -n 1 2>/dev/null | grep -E "Temperature|Fan"
```

## Notes

- Most commands work without elevated permissions.
- `sudo` required for some detailed hardware queries.
- `system_profiler` provides extensive hardware info but can be slow; use specific data types (e.g., `SPHardwareDataType`).
- For continuous monitoring, use `top -l 0` or `vm_stat 1` (they run indefinitely, use with timeout).
- Battery commands only work on MacBooks.
