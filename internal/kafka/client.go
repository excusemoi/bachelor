package kafka

import (
	"context"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/spf13/viper"
	"sync"
	"time"
)

type Client struct {
	producer *producer
	consumer *consumer
	ctx      context.Context
	metrics  *Metrics
}

func (c *Client) Init(config *viper.Viper) {
	c.ctx = context.Background()
	prefix := config.GetString("kafka.metrics.prefix")

	c.metrics = &Metrics{
		latency: &Metric{m: promauto.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_kafka_messages_latency",
			Help: "time duration between message appeared in kafka topic and obtained by component",
		}),
		},
		lag: &Metric{m: promauto.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_kafka_lag",
			Help: "amount of messages for consumer group which are waiting for processing",
		})},
		inputMessagesPerSec: &Metric{m: promauto.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_kafka_messages_input_per_second",
			Help: "amount of input messages",
		})},
		outputMessagesPerSec: &Metric{m: promauto.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_kafka_messages_output_per_second",
			Help: "amount of output messages",
		})},
	}
	c.consumer = &consumer{}
	c.producer = &producer{}
	c.producer.Init(c.ctx, config)
	c.consumer.Init(c.ctx, config)

	c.metrics.latency.count = config.GetBool("kafka.metrics.latency")
	c.metrics.lag.count = config.GetBool("kafka.metrics.lag")
	c.metrics.inputMessagesPerSec.duration = config.GetDuration("kafka.metrics.inputMessagesPerSec.duration")
	c.metrics.outputMessagesPerSec.duration = config.GetDuration("kafka.metrics.outputMessagesPerSec.duration")
}

func (c *Client) HandleEvents(handler func([]byte) ([]byte, error)) {
	wg := sync.WaitGroup{}

	if c.consumer != nil {
		if c.metrics.inputMessagesPerSec.duration != 0 {
			go c.metrics.inputMessagesPerSec.Observe(c.metrics.inputMessagesPerSec.messagePerSecondHandler)
		}
		if c.metrics.outputMessagesPerSec.duration != 0 {
			go c.metrics.outputMessagesPerSec.Observe(c.metrics.outputMessagesPerSec.messagePerSecondHandler)
		}
		wg.Add(len(c.consumer.group))
		for i := range c.consumer.group {
			go func(ind int) {
				for {
					message, err := c.consumer.group[ind].FetchMessage(c.ctx)
					if err != nil {
						wg.Done()
						break
					}

					fmt.Println("message: ", string(message.Value))

					if c.metrics.lag.count {
						c.metrics.lag.m.Set(float64(c.GetConsumerGroupLag()))
					}
					if c.metrics.latency.count {
						c.metrics.latency.m.Set(float64(time.Now().Sub(message.Time)))
					}
					if c.metrics.inputMessagesPerSec.duration != 0 {
						c.metrics.inputMessagesPerSec.value++
					}

					if handler != nil {
						go func() {
							newMessage, err := handler(message.Value)
							c.consumer.group[ind].CommitMessages(c.ctx, message)
							if err == nil {
								c.producer.produce(newMessage)
								if c.metrics.outputMessagesPerSec.duration != 0 {
									c.metrics.outputMessagesPerSec.value++
								}
							}
						}()
					}
				}
			}(i)
		}
		wg.Wait()
	}
}

func (c *Client) Produce(message []byte) error {
	return c.producer.produce(message)
}

func (c *Client) SetProducerTopics(topics []string) {
	c.producer.setTopics(topics)
}

func (c *Client) SetProducerBootstrapServers(address []string) {
	c.producer.setBootstrapServers(address)
}

func (c *Client) ProduceToTopic(message []byte, topic string) error {
	return c.producer.produceToTopic(message, topic)
}

func (c *Client) GetConsumerTopics() []string {
	return c.consumer.getTopics()
}

func (c *Client) GetConsumerGroupLag() int64 {
	lag := int64(0)
	if c.consumer != nil {
		for i := range c.consumer.group {
			lag += c.consumer.group[i].Lag()
		}
	}
	return lag
}
