package utils

import (
	"fmt"
	"time"
)

// FormatTime formats a time.Time into a string for Discord.
func FormatTime(t *time.Time, strTime string) string {
	if t == nil {
		// Use the string time if the time is nil
		return strTime
	}

	return fmt.Sprintf("<t:%d:f>", t.Unix())
}
