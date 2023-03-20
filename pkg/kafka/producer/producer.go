package producer

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer
}

func New(bs string) (*Producer, error) {
	var (
		producer *kafka.Producer
		err      error
	)

	if producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":        bs,
		"allow.auto.create.topics": true,
		"acks":                     "all",
	}); err != nil {
		return nil, err
	}
	return &Producer{producer: producer}, err
}

func (p *Producer) Run() {
	for e := range p.producer.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			{
				if ev.TopicPartition.Error != nil {

				} else {
					fmt.Printf("Successfully produced record to topic %s partition [%d] @ offset %v\n",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset,
					)
				}
			}
		}
	}
}

func (p *Producer) Produce(m *kafka.Message) error {
	return p.producer.Produce(m, nil)
}
