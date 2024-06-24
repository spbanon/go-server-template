package rabbitmq

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	amqp "github.com/rabbitmq/amqp091-go"
)

type FileReceiver struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func (l *FileReceiver) Close() {
	l.conn.Close()
	l.channel.Close()
}

func (l *FileReceiver) Receive(function func([]byte, map[string]any, string, *sync.WaitGroup) error) error {

	msgs, err := l.channel.Consume(
		l.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)

	if err != nil {
		return err
	}

	forever := make(chan os.Signal)
	signal.Notify(forever, os.Interrupt, syscall.SIGTERM)

	go func() {
		var wg sync.WaitGroup
		for d := range msgs {
			wg.Add(1)
			l.ack(d.DeliveryTag)
			go function(d.Body, d.Headers, d.RoutingKey, &wg)
		}
		wg.Wait()
	}()

	<-forever

	return nil

}

func (l *FileReceiver) ack(deliveryTag uint64) {
	err := l.channel.Ack(deliveryTag, false)
	if err != nil {
		// handle error
	}
}

func NewFileReceiver(host string, port string, username string, password string, queueName string, routingKeys []string) (*FileReceiver, error) {
	p := new(FileReceiver)

	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", username, password, host, port))
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare("tasks", "direct", true, false, false, false, map[string]any{"x-max-priority": 10})

	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		return nil, err
	}

	for _, rk := range routingKeys {

		err = ch.QueueBind(
			queue.Name, // queue name
			rk,         // routing key
			"tasks",    // exchange
			false,
			nil)

		if err != nil {
			return nil, err
		}
	}

	p.conn = conn
	p.channel = ch
	p.queueName = queue.Name

	return p, nil
}
