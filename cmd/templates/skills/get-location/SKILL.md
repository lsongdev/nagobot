---
name: get-location
description: Get the user's current geographic location.
---
# Get Location

Determine the user's current geographic location using IP geolocation or macOS CoreLocation.

## Method 1: IP-Based Geolocation (recommended, works everywhere)

Quick location lookup:
```
exec: curl -s "ipinfo.io/json"
```

This returns JSON with: `ip`, `city`, `region`, `country`, `loc` (lat,lng), `org`, `timezone`.

Get specific fields only:
```
exec: curl -s "ipinfo.io/city"
```
```
exec: curl -s "ipinfo.io/loc"
```
```
exec: curl -s "ipinfo.io/timezone"
```

## Method 2: macOS CoreLocation (precise, requires permission)

Use this when the user needs precise GPS-level location (macOS only):

```
exec: osascript -e '
use framework "CoreLocation"
use scripting additions

set locationManager to current application's CLLocationManager's alloc()'s init()
locationManager's requestWhenInUseAuthorization()
locationManager's startUpdatingLocation()

delay 3

set theLocation to locationManager's location()
if theLocation is missing value then
    return "Error: Location not available. Ensure Location Services is enabled in System Settings > Privacy & Security > Location Services."
end if

set lat to theLocation's coordinate()'s latitude() as text
set lng to theLocation's coordinate()'s longitude() as text
return "Latitude: " & lat & ", Longitude: " & lng
'
```

## Method 3: macOS Wi-Fi based (approximate)

```
exec: osascript -e 'do shell script "curl -s ipinfo.io/json"'
```

## Notes

- Method 1 is the most reliable and works on all platforms. Accuracy is city-level.
- Method 2 requires macOS Location Services to be enabled and the terminal app to have location permission.
- Combine with `get-weather` skill: get location first, then query weather for that location.
- IP geolocation may reflect VPN location if the user is connected to a VPN.
