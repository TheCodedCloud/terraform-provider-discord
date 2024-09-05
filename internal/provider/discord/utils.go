package discord

// IsSnowflake returns true if the string is a valid Discord snowflake.
func IsSnowflake(s string) bool {
	if len(s) != 18 {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}
