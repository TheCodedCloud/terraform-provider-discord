package utils

import (
	"context"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const SNOWFLAKE_MIN = 4194304 // 2^22

// IsSnowflake returns true if the string is a valid Discord snowflake.
// There are only two requirements for a valid snowflake:
// 1. It must only contain numbers.
// 2. It must be larger than SNOWFLAKE_MIN.
func IsSnowflake(s string) bool {
	// 1. Ensure that the string only contains numbers
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	// 2. Ensure that the snowflake is larger than 4194304
	snowflake, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return false
	}

	return snowflake > SNOWFLAKE_MIN
}

// User returns the current user.
func User(ctx context.Context, s *discordgo.Session) (*discordgo.User, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		result, err := s.User("@me")
		if err != nil {
			return nil, err
		}

		return result, nil
	}
}
