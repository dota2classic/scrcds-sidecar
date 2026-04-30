package rabbit

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	d2cutils "github.com/dota2classic/d2c-go-models/util"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	amqpURL  string
}

var client *Publisher

func InitRabbitPublisher() {
	host := os.Getenv("RABBITMQ_HOST")
	port := d2cutils.GetEnvInt("RABBITMQ_PORT", 5672)

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

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		log.Fatal(err)
	}

	client = &Publisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
		amqpURL:  amqpURL,
	}

	log.Println("RabbitMQ publisher initialized")
}

// reconnect tries to establish a new connection, retrying up to maxAttempts
// times before giving up. It only replaces client.conn/channel on full success,
// so the old (broken) pointers remain valid for the nil-check in publishWithRetry.
func reconnect() error {
	const maxAttempts = 10
	const delay = 2 * time.Second

	for i := 1; i <= maxAttempts; i++ {
		conn, err := amqp.Dial(client.amqpURL)
		if err != nil {
			log.Printf("rabbitmq reconnect dial attempt %d/%d failed: %v", i, maxAttempts, err)
			time.Sleep(delay)
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			conn.Close()
			log.Printf("rabbitmq reconnect channel attempt %d/%d failed: %v", i, maxAttempts, err)
			time.Sleep(delay)
			continue
		}

		err = ch.ExchangeDeclare(client.exchange, "topic", true, false, false, false, nil)
		if err != nil {
			ch.Close()
			conn.Close()
			log.Printf("rabbitmq reconnect exchange declare attempt %d/%d failed: %v", i, maxAttempts, err)
			time.Sleep(delay)
			continue
		}

		// Only swap after full success so client.channel is never nil mid-flight.
		if client.channel != nil {
			client.channel.Close()
		}
		if client.conn != nil && !client.conn.IsClosed() {
			client.conn.Close()
		}
		client.conn = conn
		client.channel = ch
		log.Println("RabbitMQ reconnected")
		return nil
	}

	return fmt.Errorf("rabbitmq reconnect failed after %d attempts", maxAttempts)
}

func isConnectionError(err error) bool {
	if errors.Is(err, amqp.ErrClosed) {
		return true
	}
	var amqpErr *amqp.Error
	// 504 = CHANNEL-ERROR, 320 = CONNECTION-FORCED
	return errors.As(err, &amqpErr) && (amqpErr.Code == 504 || amqpErr.Code == 320)
}

func publishWithRetry[T any](event *T, routingKey string, retries int) error {
	message, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("can't serialize payload: %w", err)
	}

	for attempt := 1; attempt <= retries; attempt++ {
		// Guard: if reconnect previously failed we may have a broken channel;
		// check liveness before publishing.
		if client.channel == nil || client.conn == nil || client.conn.IsClosed() {
			log.Printf("rabbitmq channel unavailable before attempt %d/%d — reconnecting...", attempt, retries)
			if reconnErr := reconnect(); reconnErr != nil {
				return fmt.Errorf("rabbitmq publish aborted, cannot reconnect: %w", reconnErr)
			}
		}

		err = client.channel.Publish(
			client.exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        message,
				Timestamp:   time.Now(),
			},
		)

		if err == nil {
			return nil
		}

		// Any error from Publish is treated as a connection-level problem and
		// triggers reconnect. Non-connection AMQP errors (e.g. wrong exchange
		// name) would also land here, but in practice those are config issues
		// that would have failed at startup. This is safer than only matching
		// known error codes and risking a missed case that drops the event.
		log.Printf("rabbitmq publish failed (attempt %d/%d): %v — reconnecting...", attempt, retries, err)
		if reconnErr := reconnect(); reconnErr != nil {
			// reconnect already retried internally; nothing more to do this round.
			log.Printf("rabbitmq reconnect failed: %v", reconnErr)
		}
	}

	return fmt.Errorf("rabbitmq publish failed after %d retries: %w", retries, err)
}

func (r *Publisher) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
