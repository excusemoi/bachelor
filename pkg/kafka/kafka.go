package kafka

import (
	"github.com/bachelor/pkg/kafka/consumer"
	"github.com/bachelor/pkg/kafka/producer"
	"github.com/bachelor/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Client struct {
	Producer *producer.Producer
	Consumer *consumer.Consumer
	Filter   func(*kafka.Message, *model.FiltrationRule) (*kafka.Message, error) //must be interface with Filter() method
}

func New(bs string, topics []string) (*Client, error) {
	var (
		err error
		c   = &Client{}
	)
	if c.Consumer, err = consumer.New(bs, topics); err != nil {
		return nil, err
	}
	if c.Producer, err = producer.New(bs); err != nil {
		return nil, err
	}
	return c, err
}

func (c *Client) Run() {
	go c.Producer.Run()
	c.Consumer.Run(c.HandleMessage)
}

func (c *Client) HandleMessage(m *kafka.Message) error {
	var err error

	//if m, err = c.Filter(m, fr); err != nil {
	//	return err
	//}

	if err = c.Producer.Produce(m); err != nil {
		return err
	}

	return err
}
