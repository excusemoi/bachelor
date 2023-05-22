package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
)

type consumer struct {
	group  []*kafka.Reader
	topics []string
	ctx    context.Context
}

func (c *consumer) Init(ctx context.Context, vp *viper.Viper) {
	c.ctx = ctx
	c.topics = vp.GetStringSlice("kafka.consumer.topics")
	c.group = make([]*kafka.Reader, len(c.topics))
	bootstrapServers := vp.GetStringSlice("kafka.bootstrapServers")
	for i := range c.topics {
		c.group[i] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:         bootstrapServers,
			GroupID:         vp.GetString("kafka.consumer.groupId"),
			Topic:           c.topics[i],
			RetentionTime:   0,
			ReadLagInterval: 0,
		},
		)
	}
}

func (c *consumer) getTopics() []string {
	return c.topics
}
