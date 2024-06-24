package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type FileSender struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchName string
	mu       sync.Mutex
}

func (l *FileSender) Close() {
	l.conn.Close()
	l.channel.Close()
}

func (l *FileSender) Send(dump any, routingKey string, headers map[string]any) error {
	filesJson, err := json.Marshal(dump)
	if err != nil {
		return err
	}

	msgEncoded := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Body:         filesJson,
		Headers:      headers,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	l.mu.Lock()
	defer l.mu.Unlock()
	err = l.channel.PublishWithContext(ctx, l.exchName, routingKey, false, false, msgEncoded)
	return err

}

func NewFileSender(host string, port string, username string, password string, exchangeName string) (*FileSender, error) {
	p := new(FileSender)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(
		exchangeName,                         // name
		"direct",                             // type
		true,                                 // durable
		false,                                // auto-deleted
		false,                                // internal
		false,                                // no-wait
		map[string]any{"x-max-priority": 10}, // arguments
	)
	if err != nil {
		return nil, err
	}

	p.conn = conn
	p.channel = ch
	p.exchName = exchangeName

	return p, nil
}
