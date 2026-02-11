---
name: get-weather
description: Get current weather and forecast for a location.
---
# Get Weather

Query current weather conditions and forecasts using wttr.in (no API key required).

## Workflow

1. **Get current weather for a city** (concise, one-line):
   ```
   exec: curl -s "wttr.in/CityName?format=%l:+%C+%t+%h+%w"
   ```

2. **Get detailed weather report** (multi-line with forecast):
   ```
   exec: curl -s "wttr.in/CityName?format=v2&lang=zh"
   ```

3. **Get weather as JSON** (for programmatic use):
   ```
   exec: curl -s "wttr.in/CityName?format=j1"
   ```

4. **Get weather for user's current location** (IP-based):
   ```
   exec: curl -s "wttr.in/?format=%l:+%C+%t+%h+%w"
   ```

## Format Placeholders

- `%l` — location
- `%C` — weather condition text
- `%t` — temperature
- `%h` — humidity
- `%w` — wind
- `%p` — precipitation (mm)
- `%P` — pressure
- `%D` — dawn time
- `%S` — sunrise
- `%s` — sunset

Custom format example (temperature + condition + wind):
```
exec: curl -s "wttr.in/Beijing?format=%l:+%C+%t+wind:+%w"
```

## Language Support

Append `&lang=xx` to get results in a specific language:
- `lang=zh` — Chinese
- `lang=en` — English
- `lang=ja` — Japanese

## Notes

- City names with spaces should be URL-encoded or use `+` (e.g., `New+York`).
- Use `~CityName` for unicode location names (e.g., `~北京`).
- wttr.in is rate-limited; avoid excessive requests.
- If `curl` is not available, use `web_fetch` tool with the same URLs.
