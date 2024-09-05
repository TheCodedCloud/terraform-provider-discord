package common

import (
	"time"
)

// StrDiscordTimestamp formats a Discord timestamp.
func StrDiscordTime(timestamp *time.Time, format string) string {
	if timestamp == nil {
		return ""
	}

	switch format {
	case "ISO8601":
		format = time.RFC3339
	case "RFC850":
		format = time.RFC850
	}

	if format == "" {
		format = time.RFC850
	}

	return timestamp.Format(format)
}

// CurrentTime returns the current time as a string.
func CurrentTime() string {
	return time.Now().Format(time.RFC850)
}
