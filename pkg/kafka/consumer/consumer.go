package consumer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"
)

type Consumer struct {
	consumer *kafka.Consumer
}

func New(bs string, topics []string) (*Consumer, error) {
	var (
		consumer *kafka.Consumer
		err      error
	)

	if consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":               bs,
		"group.id":                        "group",
		"auto.offset.reset":               "earliest",
		"go.application.rebalance.enable": true}); err != nil {
		return nil, err
	}

	if err = consumer.SubscribeTopics(topics, nil); err != nil {
		return nil, err
	}

	return &Consumer{consumer: consumer}, err
}

func (c *Consumer) Run(messageHandler func(message *kafka.Message) error) {
	run := true

	for run {
		ev := c.consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			go func() {
				messageHandler(e)
				c.consumer.Commit()
			}()
			fmt.Printf("%% Message on %s:\n%s\n",
				e.TopicPartition, string(e.Value))

		case kafka.PartitionEOF:
			fmt.Printf("%% Reached %v\n", e)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
			fmt.Printf("Ignored %v\n", e)
		}
	}
}
