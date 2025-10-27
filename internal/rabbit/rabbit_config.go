package rabbit

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"sidecar/internal/util"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

var client *Publisher

func InitRabbitPublisher() {
	host := os.Getenv("RABBITMQ_HOST")
	port := util.GetEnvInt("RABBITMQ_PORT", 5672)

	username := os.Getenv("RABBITMQ_USER")
	password := os.Getenv("RABBITMQ_PASSWORD")

	exchange := "app.events"

	amqpURL := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)

	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		log.Fatal(err)
	}

	client = &Publisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}
}

func publishWithRetry[T any](event T, retries int) error {
	var err error

	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("can't serialize payload: %w", err)
	}

	for attempt := 1; attempt <= retries; attempt++ {
		err = client.channel.Publish(
			client.exchange, // exchange
			"",              // routing key (empty for fanout)
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        message,
				Timestamp:   time.Now(),
			},
		)

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

// Close the connection and channel
func (r *Publisher) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
