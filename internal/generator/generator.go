package generator

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Generator struct {
	producer *kafka.Writer
	ctx      context.Context
	topics   []string
}

func New() *Generator {
	return &Generator{producer: &kafka.Writer{}, topics: []string{}}
}

func (g *Generator) Init(path, name string) error {
	var (
		vp  = viper.New()
		err error
	)

	vp.SetConfigName(name)
	vp.SetConfigType("yaml")
	vp.AddConfigPath(path)
	if err := vp.ReadInConfig(); err != nil {
		return err
	}
	g.ctx = context.Background()
	g.topics = vp.GetStringSlice("kafka.producer.topics")
	g.producer = &kafka.Writer{
		Addr:                   kafka.TCP(vp.GetStringSlice("kafka.bootstrapServers")...),
		Async:                  true,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true}
	return err
}

func (g *Generator) Run() error {

	var (
		dataFile *os.File
		data     []byte
		messages []map[string]interface{}
		err      error
	)

	if dataFile, err = os.Open(filepath.Join("..", "data", "data.json")); err != nil {
		return err
	}

	if data, err = ioutil.ReadAll(dataFile); err != nil {
		return err
	}

	if err = json.Unmarshal(data, &messages); err != nil {
		return err
	}

	for {
		for {
			for i := range messages {
				message, _ := json.Marshal(messages[i])
				msgs := make([]kafka.Message, len(g.topics))
				for j := range g.topics {
					msgs[j] = kafka.Message{
						Topic:     g.topics[j],
						Value:     message,
						Partition: 0}
				}
				go g.producer.WriteMessages(g.ctx, msgs...)
			}
		}
	}
}
