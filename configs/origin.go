package configs

import (
	"os"
	"strings"
)

func LoadAllowedOrigins() map[string]bool {
	origins := make(map[string]bool)
	envOrigins := os.Getenv("ALLOWED_ORIGINS")
	for _, origin := range strings.Split(envOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins[origin] = true
		}
	}
	return origins
}
