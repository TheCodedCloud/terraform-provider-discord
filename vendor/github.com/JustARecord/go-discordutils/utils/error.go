package utils

import "strings"

// NotFoundError is an error type for when a resource is not found.
func NotFoundError(err error) bool {
	return err != nil && strings.HasPrefix(err.Error(), "HTTP 404 Not Found")
}
