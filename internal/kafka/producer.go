package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type producer struct {
	writer *kafka.Writer
	topics []string
	ctx    context.Context
}

func (p *producer) Init(ctx context.Context, vp *viper.Viper) {
	p.ctx = ctx
	p.topics = vp.GetStringSlice("kafka.producer.topics")
	p.writer = &kafka.Writer{
		Addr:                   kafka.TCP(vp.GetStringSlice("kafka.bootstrapServers")...),
		Async:                  true,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

func (p *producer) produce(message []byte) error {
	messages := make([]kafka.Message, len(p.topics))
	for i := range p.topics {
		messages[i] = kafka.Message{
			Topic:     p.topics[i],
			Value:     message,
			Partition: 0,
		}
	}
	return p.writer.WriteMessages(p.ctx, messages...)
}

func (p *producer) produceToTopic(message []byte, topic string) error {
	m := []kafka.Message{{Topic: topic, Value: message, Partition: 0}}
	return p.writer.WriteMessages(p.ctx, m...)
}

func (p *producer) setTopics(topics []string) {
	p.topics = topics
}

func (p *producer) setBootstrapServers(address []string) {
	p.writer.Addr = kafka.TCP(address...)
}
