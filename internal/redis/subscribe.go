package redis

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Subscribe subscribes to a Redis channel and calls handler for each message.
// It runs the receive loop in a separate goroutine.
func Subscribe[In any, Out any](
	ctx context.Context,
	client *redis.Client,
	channel string,
	handler func(*In) (*Out, error),
) {
	go func() {
		backoff := time.Second

		for {
			select {
			case <-ctx.Done():
				log.Printf("[RedisBus] Subscription to %s cancelled", channel)
				return
			default:
			}

			pubsub := client.Subscribe(ctx, channel)

			// Wait for the subscription to establish
			if _, err := pubsub.Receive(ctx); err != nil {
				log.Printf("[RedisBus] Subscribe error: %v", err)
				pubsub.Close()
				time.Sleep(backoff)
				backoff = min(backoff*2, 30*time.Second)
				continue
			}

			log.Printf("[RedisBus] Subscribed to %s", channel)
			backoff = time.Second // reset backoff on success

			ch := pubsub.Channel()

			// Listen loop
			for {
				select {
				case msg, ok := <-ch:
					if !ok {
						log.Printf("[RedisBus] Channel %s closed, reconnecting...", channel)
						pubsub.Close()
						time.Sleep(backoff)
						backoff = min(backoff*2, 5*time.Second)
						break
					}

					log.Printf("[RedisBus] Channel %s received: %s", channel, msg.Payload)

					var payload ChannelEvent[In]
					if err := json.Unmarshal([]byte(msg.Payload), &payload); err != nil {
						log.Printf("[RedisBus] Invalid message on %s: %v", channel, err)
						continue
					}

					log.Printf("[RedisBus] Channel %s unmarshalled %v", channel, payload)

					// Safely call handler
					func() {
						defer func() {
							if r := recover(); r != nil {
								log.Printf("[RedisBus] Panic in handler: %v", r)
							}
						}()
						log.Printf("Handling message in channel %s", channel)

						res, err := handler(&payload.Data)
						if err != nil {
							log.Printf("[RedisSubscribe] Handler error: %v", err)
						}

						if res != nil {
							//
							log.Printf("[RedisSubscribe] Handler returned %v", res)
						}

						replyChannel := channel + ".reply"

						response := ChannelEvent[Out]{
							Id:      payload.Id,
							Data:    *res,
							Pattern: payload.Pattern,
						}

						bt, err := json.Marshal(response)

						log.Printf("[RedisSubscribe] Publishing message to %s %v", channel, response)
						client.Publish(ctx, replyChannel, bt)

					}()

				case <-ctx.Done():
					pubsub.Close()
					log.Printf("[RedisBus] Unsubscribed from %s", channel)
					return
				}
			}
		}
	}()
}
