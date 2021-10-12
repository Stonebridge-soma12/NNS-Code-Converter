package messageQ

import (
	"fmt"
	"github.com/labstack/gommon/random"
	"github.com/streadway/amqp"
	"time"
)

type ResponseMsg struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type MessageQ struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

const (
	defaultHeartbeat = 10 * time.Second
	defaultLocale    = "en_US"
)

func CreateConnection(id, pw, host, vhost string) (*MessageQ, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", id, pw, host)
	config := amqp.Config{
		Heartbeat: defaultHeartbeat,
		Locale:    defaultLocale,
		Vhost:     vhost,
	}
	conn, err := amqp.DialConfig(url, config)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	res := &MessageQ{connection: conn, channel: ch}

	return res, nil
}

func (m *MessageQ) Publish(data []byte) error {
	requestId := random.String(32)
	publish := amqp.Publishing{
		DeliveryMode:  amqp.Persistent,
		ContentType:   "application/json",
		CorrelationId: requestId,
		Body:          data,
	}

	err := m.channel.Publish("", "Request", false, false, publish)
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageQ) Close() error {
	err := m.channel.Close()
	if err != nil {
		return err
	}

	err = m.connection.Close()
	if err != nil {
		return err
	}

	return nil
}
