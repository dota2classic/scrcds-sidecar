package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	d2cmodels "github.com/dota2classic/d2c-go-models/models"
	"github.com/redis/go-redis/v9"
)

func PublishLiveMatch(evt *d2cmodels.LiveMatchUpdateEvent) {
	err := publishWithRetry("LiveMatchUpdateEvent", evt, 3)
	if err != nil {
		fmt.Printf("There was an issue publishing event: %s\n", err)
	}
}

func PublishPlayerConnectedEvent(evt *d2cmodels.PlayerConnectedEvent) {
	err := publishWithRetry("PlayerConnectedEvent", evt, 3)
	if err != nil {
		fmt.Printf("There was an issue publishing event: %s\n", err)
	}
}

// publishWithRetry publishes a message with automatic retry logic.
func publishWithRetry[T any](channel string, event *T, retries int) error {
	var err error

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("can't serialize payload: %w", err)
	}

	for attempt := 1; attempt <= retries; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Publish(ctx, channel, message).Err()
		cancel()

		if err == nil {
			return nil
		}

		// Check if it's a transient error (network, closed conn, etc.)
		if errors.Is(err, redis.ErrClosed) || errors.Is(err, context.DeadlineExceeded) {
			fmt.Printf("⚠️  Redis publish failed (attempt %d/%d): %v — retrying...\n", attempt, retries, err)
			time.Sleep(time.Duration(attempt) * 200 * time.Millisecond)
			continue
		}

		// Non-retryable error
		return fmt.Errorf("redis publish failed: %w", err)
	}

	return fmt.Errorf("redis publish failed after %d retries: %w", retries, err)
}
