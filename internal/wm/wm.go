package wm

import (
	"os"
	"strings"
)

func Detect() string {
	env := os.Getenv("XDG_CURRENT_DESKTOP")
	if env != "" {
		return strings.ToLower(env)
	}
	session := os.Getenv("DESKTOP_SESSION")
	if session != "" {
		return strings.ToLower(session)
	}
	return "dwl" // Default fallback for your setup
}
