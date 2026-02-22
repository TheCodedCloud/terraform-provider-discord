package utils

import (
	"context"
	"log/slog"
	"time"
)

// Sleep sleeps for the given duration, but can be interrupted by context
func Sleep(ctx context.Context, duration time.Duration) {
	slog.Info("Sleeping", "interval", duration)
	select {
	case <-ctx.Done():
		slog.Info("Sleep interrupted", "duration", duration)
	case <-time.After(duration):
	}
}
