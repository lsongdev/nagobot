package provider

import "strings"

func normalizeSDKBaseURL(raw, defaultBase string, endpointSuffixes ...string) string {
	base := strings.TrimSpace(raw)
	if base == "" {
		return defaultBase
	}

	base = strings.TrimRight(base, "/")
	for _, suffix := range endpointSuffixes {
		s := strings.TrimRight(strings.TrimSpace(suffix), "/")
		if s == "" {
			continue
		}
		if strings.HasSuffix(base, s) {
			base = strings.TrimSuffix(base, s)
			base = strings.TrimRight(base, "/")
			break
		}
	}

	if base == "" {
		return defaultBase
	}
	return base
}
