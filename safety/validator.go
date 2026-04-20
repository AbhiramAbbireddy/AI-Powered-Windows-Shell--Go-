package safety

import (
	"strings"
)

var dangerousPatterns = []string{
	"del *",
	"erase *",
	"format",
	"shutdown",
	"rd /s",
	"rmdir /s",
}

func IsDangerous(command string) (bool, string) {
	normalized := strings.ToLower(strings.TrimSpace(command))
	for _, pattern := range dangerousPatterns {
		if strings.Contains(normalized, pattern) {
			return true, "command can delete data or affect the system"
		}
	}

	if strings.Contains(normalized, "/q") && (strings.HasPrefix(normalized, "del ") || strings.HasPrefix(normalized, "rd ")) {
		return true, "quiet destructive command detected"
	}

	return false, ""
}
