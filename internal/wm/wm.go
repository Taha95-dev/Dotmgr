package wm

import (
	"os"
	"strings"
)

func Detect() string {
	env := os.Getenv("XDG_CURRENT_DESKTOP")
	if env != "" {
		return normalize(env)
	}
	session := os.Getenv("DESKTOP_SESSION")
	if session != "" {
		return normalize(session)
	}
	return "dwl" // Default fallback for your setup
}

func normalize(raw string) string {
	wm := strings.ToLower(strings.TrimSpace(raw))
	if idx := strings.Index(wm, ":"); idx != -1 {
		wm = wm[:idx]
	}
	return wm
}

