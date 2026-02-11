---
name: docker-manage
description: Docker container and image management (ps, logs, build, compose).
tags: [docker, dev, cross-platform]
---
# Docker Management

Cross-platform Docker operations. Requires Docker Engine or Docker Desktop.

## Containers

List running containers:
```
exec: docker ps
```

List all containers (including stopped):
```
exec: docker ps -a
```

Start / stop / restart:
```
exec: docker start CONTAINER
```
```
exec: docker stop CONTAINER
```
```
exec: docker restart CONTAINER
```

Remove container:
```
exec: docker rm CONTAINER
```

Remove all stopped containers:
```
exec: docker container prune -f
```

## Logs

Tail logs:
```
exec: docker logs --tail 100 CONTAINER
```

Follow logs (use with timeout):
```
exec: docker logs -f --tail 50 CONTAINER
```

Logs since timestamp:
```
exec: docker logs --since "2026-02-10T00:00:00" CONTAINER
```

## Exec Into Container

Interactive shell:
```
exec: docker exec -it CONTAINER /bin/sh
```

Run a command:
```
exec: docker exec CONTAINER cat /etc/os-release
```

## Images

List images:
```
exec: docker images
```

Pull image:
```
exec: docker pull IMAGE:TAG
```

Build image:
```
exec: docker build -t IMAGE_NAME:TAG -f Dockerfile .
```

Remove image:
```
exec: docker rmi IMAGE
```

Remove unused images:
```
exec: docker image prune -f
```

## Docker Compose

Start services:
```
exec: docker compose -f /path/to/docker-compose.yml up -d
```

Stop services:
```
exec: docker compose -f /path/to/docker-compose.yml down
```

View service logs:
```
exec: docker compose -f /path/to/docker-compose.yml logs --tail 50 SERVICE
```

List services:
```
exec: docker compose -f /path/to/docker-compose.yml ps
```

Rebuild and restart:
```
exec: docker compose -f /path/to/docker-compose.yml up -d --build SERVICE
```

## Resource Usage

Container resource stats:
```
exec: docker stats --no-stream
```

Disk usage:
```
exec: docker system df
```

## Networking

List networks:
```
exec: docker network ls
```

Inspect network:
```
exec: docker network inspect NETWORK_NAME
```

## Volumes

List volumes:
```
exec: docker volume ls
```

Remove unused volumes:
```
exec: docker volume prune -f
```

## Full Cleanup

Remove all unused data (containers, images, networks, volumes):
```
exec: docker system prune -a --volumes -f
```

## Inspect

Container details:
```
exec: docker inspect CONTAINER | head -80
```

Container IP address:
```
exec: docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' CONTAINER
```

## Notes

- Works on macOS, Linux, and Windows (Docker Desktop or WSL).
- `CONTAINER` can be container name or ID.
- Use `docker compose` (v2) instead of `docker-compose` (v1, deprecated).
- `docker stats --no-stream` gives a one-time snapshot; without `--no-stream` it runs continuously.
- Some commands require the Docker daemon to be running.
