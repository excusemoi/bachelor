package actor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
	"strings"
)

type Actor struct {
	components.AbstractComponent[model.ActionRule]
	TopicDeterminer func(string) string
}

func (f *Actor) Init(path string, name string, determiner func(string) string) error {
	f.TopicDeterminer = determiner
	return f.AbstractComponent.Init(path, name)
}

func (f *Actor) Handle(message []byte) ([]byte, error) {
	var (
		m   = map[string]interface{}{}
		res interface{}
		in  bool
	)
	for _, rule := range f.Rules {
		if err := json.Unmarshal(message, &m); err != nil {
			return nil, errors.New("actor: can't unmarshall message")
		}
		if res, in = m[rule.Key]; !in {
			return nil, errors.New(fmt.Sprintf("actor: no field %s in message %s", rule.Key, string(message)))
		}
		f.Kafka.SetProducerBootstrapServers(strings.Split(rule.BootstrapServers, ","))
		f.Kafka.SetProducerTopics([]string{f.TopicDeterminer(res.(string))})
	}
	return json.Marshal(&m)
}

func (f *Actor) Run() {
	log.Println("actor running")
	go f.Observe()
	f.Kafka.HandleEvents(f.Handle)
}
