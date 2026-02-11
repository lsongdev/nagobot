---
name: network-diagnose
description: Network diagnostics (ping, traceroute, dig, port scan, bandwidth).
tags: [network, diagnostics, cross-platform]
---
# Network Diagnostics

Diagnose network issues using standard cross-platform tools.

## Connectivity Check

Ping host:
```
exec: ping -c 4 example.com
```

Quick reachability test:
```
exec: ping -c 1 -W 3 8.8.8.8 && echo "Online" || echo "Offline"
```

## DNS Lookup

Basic lookup:
```
exec: dig example.com +short
```

Full record:
```
exec: dig example.com
```

Specific record type:
```
exec: dig example.com MX +short
```
```
exec: dig example.com AAAA +short
```
```
exec: dig example.com TXT +short
```

Reverse DNS:
```
exec: dig -x 8.8.8.8 +short
```

Using specific DNS server:
```
exec: dig @1.1.1.1 example.com +short
```

Fallback (if `dig` not available):
```
exec: nslookup example.com
```
```
exec: host example.com
```

## Route Tracing

Traceroute:
```
exec: traceroute -m 20 example.com
```

On Linux (may need `traceroute` package):
```
exec: traceroute example.com
```

macOS alternative:
```
exec: traceroute -q 1 example.com
```

## Port Checking

Check if a port is open:
```
exec: nc -zv -w 3 example.com 443 2>&1
```

Scan port range:
```
exec: for port in 22 80 443 3306 5432 6379 8080; do nc -zv -w 2 example.com $port 2>&1; done
```

Check local listening ports:
```
exec: ss -tlnp 2>/dev/null || netstat -tlnp 2>/dev/null || lsof -i -P -n | grep LISTEN
```

## Public IP

```
exec: curl -s ifconfig.me
```

Detailed IP info:
```
exec: curl -s ipinfo.io/json
```

## Network Interfaces

List interfaces:
```
exec: ip addr 2>/dev/null || ifconfig
```

Active interface and gateway:
```
exec: ip route 2>/dev/null || netstat -rn | head -10
```

## Bandwidth Test (Simple)

Download speed:
```
exec: curl -s -o /dev/null -w "Speed: %{speed_download} bytes/sec\nTime: %{time_total}s\nSize: %{size_download} bytes\n" https://speed.cloudflare.com/__down?bytes=10000000
```

## Whois

Domain info:
```
exec: whois example.com | head -40
```

## SSL/TLS Certificate

Check certificate:
```
exec: echo | openssl s_client -connect example.com:443 -servername example.com 2>/dev/null | openssl x509 -noout -dates -subject -issuer
```

Certificate expiry only:
```
exec: echo | openssl s_client -connect example.com:443 -servername example.com 2>/dev/null | openssl x509 -noout -enddate
```

## HTTP Headers

```
exec: curl -sI https://example.com
```

## Active Connections

```
exec: ss -tnp 2>/dev/null || netstat -tn
```

Connection count by state:
```
exec: ss -tn 2>/dev/null | awk '{print $1}' | sort | uniq -c | sort -rn || netstat -tn | awk '{print $6}' | sort | uniq -c | sort -rn
```

## Notes

- `ping`, `curl`, `openssl` are available on nearly all systems.
- `dig` may need `dnsutils` (Linux) or `bind-tools` package.
- `traceroute` may need `traceroute` package on minimal Linux installs.
- `ss` is the modern replacement for `netstat` on Linux; `netstat` works on macOS.
- Some commands may require `sudo` (e.g., `tcpdump`, certain `nmap` scans).
- The `nc` (netcat) command syntax varies slightly between macOS and Linux.
