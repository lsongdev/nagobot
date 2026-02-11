---
name: archive
description: Compress and extract archives (tar, zip, gzip, 7z).
tags: [archive, compress, files, cross-platform]
---
# Archive & Compression

Create, extract, and manage compressed archives. Cross-platform (macOS, Linux).

## tar + gzip (.tar.gz / .tgz)

Create:
```
exec: tar -czf /path/to/archive.tar.gz -C /path/to/parent directory_name
```

Extract:
```
exec: tar -xzf /path/to/archive.tar.gz -C /path/to/destination
```

List contents:
```
exec: tar -tzf /path/to/archive.tar.gz
```

## tar + bzip2 (.tar.bz2)

Create:
```
exec: tar -cjf /path/to/archive.tar.bz2 -C /path/to/parent directory_name
```

Extract:
```
exec: tar -xjf /path/to/archive.tar.bz2 -C /path/to/destination
```

## tar + xz (.tar.xz) — best compression ratio

Create:
```
exec: tar -cJf /path/to/archive.tar.xz -C /path/to/parent directory_name
```

Extract:
```
exec: tar -xJf /path/to/archive.tar.xz -C /path/to/destination
```

## zip (.zip)

Create:
```
exec: zip -r /path/to/archive.zip /path/to/directory
```

Extract:
```
exec: unzip /path/to/archive.zip -d /path/to/destination
```

List contents:
```
exec: unzip -l /path/to/archive.zip
```

Add files to existing zip:
```
exec: zip -g /path/to/archive.zip /path/to/newfile.txt
```

Create with password:
```
exec: zip -e -r /path/to/secure.zip /path/to/directory
```

## gzip (single file)

Compress:
```
exec: gzip /path/to/file.txt
```

Decompress:
```
exec: gunzip /path/to/file.txt.gz
```

Keep original:
```
exec: gzip -k /path/to/file.txt
```

## 7-Zip (.7z) — requires p7zip

Create:
```
exec: 7z a /path/to/archive.7z /path/to/directory
```

Extract:
```
exec: 7z x /path/to/archive.7z -o/path/to/destination
```

List:
```
exec: 7z l /path/to/archive.7z
```

## Inspect Any Archive

File type detection:
```
exec: file /path/to/archive
```

## Exclude Patterns

Tar with exclusions:
```
exec: tar -czf /path/to/archive.tar.gz --exclude='*.log' --exclude='node_modules' -C /path/to/parent directory_name
```

Zip with exclusions:
```
exec: zip -r /path/to/archive.zip /path/to/directory -x "*.log" "*/node_modules/*"
```

## Split Large Archives

Create split archive (100MB parts):
```
exec: tar -czf - /path/to/large_dir | split -b 100m - /path/to/archive.tar.gz.part
```

Reassemble and extract:
```
exec: cat /path/to/archive.tar.gz.part* | tar -xzf - -C /path/to/destination
```

## Compare Archive Contents

```
exec: diff <(tar -tzf archive1.tar.gz | sort) <(tar -tzf archive2.tar.gz | sort)
```

## Notes

- `tar`, `gzip`, `zip`/`unzip` are pre-installed on macOS and most Linux distros.
- `xz` offers the best compression ratio but is slower.
- `7z` needs `p7zip` or `p7zip-full` (install via `brew install p7zip` or `apt install p7zip-full`).
- Use `-C` with tar to control the base directory in the archive.
- For Windows compatibility, `.zip` is the safest format.
