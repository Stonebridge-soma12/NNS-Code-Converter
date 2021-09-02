package MessageQ

import (
	"fmt"
	"github.com/streadway/amqp"
)

type ResponseMsg struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

func CreateConnection(id, pw, host string) (*amqp.Connection, error) {
	fmt.Println(pw)
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", id, pw, host)
	fmt.Println(url)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
