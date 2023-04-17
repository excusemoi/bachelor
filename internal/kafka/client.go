package kafka

import (
	"context"
	"github.com/spf13/viper"
	"sync"
)

type Client struct {
	producer *producer
	consumer *consumer
	ctx      context.Context
}

func (c *Client) Init(config *viper.Viper) {
	c.ctx = context.Background()
	c.consumer = &consumer{}
	c.producer = &producer{}
	c.producer.Init(c.ctx, config)
	c.consumer.Init(c.ctx, config)
}

func (c *Client) HandleEvents(handler func([]byte) error) {
	wg := sync.WaitGroup{}
	wg.Add(len(c.consumer.group))
	for i := range c.consumer.group {
		go func(ind int) {
			for {
				message, err := c.consumer.group[ind].FetchMessage(c.ctx)
				if err != nil {
					//log.ErrorSimple(err, "kafka", "consumer.handleMessage") TODO log error while fetching message
					wg.Done()
					break
				}
				//log.InfoSimple("kafka", "consumer.handleMessage", TODO log obtained message
				//	fmt.Sprintf("message at topic/partition/offset %v/%v/%v: %s = %s\n",
				//		message.Topic,
				//		message.Partition,
				//		message.Offset, string(message.Key),
				//		string(message.Value)))
				if handler != nil {
					if err = handler(message.Value); err != nil {
						//TODO log cant handle message
					} else {
						c.consumer.group[ind].CommitMessages(c.ctx, message)
					}
				}
			}
		}(i)
	}
	wg.Wait()
}

func (c *Client) Produce(message []byte) error {
	return c.producer.produce(message)
}

func (c *Client) SetProducerTopics(topics []string) {
	c.producer.setTopics(topics)
}

func (c *Client) GetConsumerTopics() []string {
	return c.consumer.getTopics()
}
