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
	fmt.Println(config.GetString("kafka.metrics.prefix"))

	if config.Get("kafka.metrics") != nil {
		c.metrics = &Metrics{
			latency:              &Metric{},
			lag:                  &Metric{},
			inputMessagesPerSec:  &Metric{},
			outputMessagesPerSec: &Metric{},
			filtrationParams:     &Metric{},
		}
		c.metrics.latency.count = config.GetBool("kafka.metrics.latency")
		if c.metrics.latency.count {
			c.metrics.latency.m = promauto.NewGauge(
				prometheus.GaugeOpts{
					Name: prefix + "_kafka_messages_latency",
					Help: "time start between message appeared in kafka topic and obtained by component",
				})
		}

		c.metrics.lag.count = config.GetBool("kafka.metrics.lag")
		if c.metrics.lag.count {
			c.metrics.lag.m = promauto.NewGauge(prometheus.GaugeOpts{
				Name: prefix + "_kafka_lag",
				Help: "amount of messages for consumer group which are waiting for processing",
			})
		}

		c.metrics.filtrationParams.count = config.GetBool("kafka.metrics.filtrationParams")
		if c.metrics.filtrationParams.count {
			c.metrics.filtrationParams.m = promauto.NewGauge(prometheus.GaugeOpts{
				Name: prefix + "_kafka_filtration_params",
				Help: "amount of filtration params",
			})
		}

		c.metrics.inputMessagesPerSec.count = config.GetBool("kafka.metrics.inputMessagesPerSec.count")
		if c.metrics.inputMessagesPerSec.count {
			c.metrics.inputMessagesPerSec.m = promauto.NewGauge(prometheus.GaugeOpts{
				Name: prefix + "_kafka_messages_input_per_second",
				Help: "amount of input messages",
			})
		}

		c.metrics.outputMessagesPerSec.count = config.GetBool("kafka.metrics.outputMessagesPerSec.count")
		if c.metrics.outputMessagesPerSec.count {
			c.metrics.outputMessagesPerSec.m = promauto.NewGauge(prometheus.GaugeOpts{
				Name: prefix + "_kafka_messages_output_per_second",
				Help: "amount of output messages",
			})
		}
	}
	c.consumer = &consumer{}
	c.producer = &producer{}
	c.producer.Init(c.ctx, config)
	c.consumer.Init(c.ctx, config)
}

func (c *Client) HandleEvents(handler func([]byte) ([]byte, error)) {
	wg := sync.WaitGroup{}
	if c.consumer != nil {
		wg.Add(len(c.consumer.group))
		if c.metrics.inputMessagesPerSec.count {
			c.metrics.inputMessagesPerSec.start = time.Now()
		}
		if c.metrics.outputMessagesPerSec.count {
			c.metrics.outputMessagesPerSec.start = time.Now()
		}
		for i := range c.consumer.group {
			go func(ind int) {
				buff := make(chan struct{}, 1000000)
				for {
					message, err := c.consumer.group[ind].FetchMessage(c.ctx)

					fmt.Println("message: ", message)

					if err != nil {
						wg.Done()
						return
					}

					if c.metrics.lag.count {
						c.metrics.lag.m.Set(float64(c.GetConsumerGroupLag()))
					}
					if c.metrics.latency.count {
						latency := float64(time.Now().Sub(message.Time).Milliseconds())
						if latency < 100. {
							c.metrics.latency.m.Set(latency)

						}
					}
					if c.metrics.inputMessagesPerSec.count {
						c.metrics.inputMessagesPerSec.value++
						c.metrics.inputMessagesPerSec.m.Set(float64(c.metrics.inputMessagesPerSec.value) /
							time.Now().Sub(c.metrics.inputMessagesPerSec.start).Seconds())
					}

					if handler != nil {
						buff <- struct{}{}
						go func() {
							newMessage, err := handler(message.Value)
							c.consumer.group[ind].CommitMessages(c.ctx, message)
							if err == nil {
								c.producer.produce(newMessage)
								if c.metrics.outputMessagesPerSec.count {
									c.metrics.outputMessagesPerSec.value++
									c.metrics.outputMessagesPerSec.m.Set(float64(c.metrics.outputMessagesPerSec.value) /
										time.Now().Sub(c.metrics.outputMessagesPerSec.start).Seconds())
								}
							}
							<-buff
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
			l := c.consumer.group[i].Stats().Lag
			if l < 30000000 {
				lag += l
			}
		}
	}
	return lag
}
