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

func (f *Actor) Init(configPath string, determiner func(string) string) error {
	if err := f.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	f.Rule = &model.ActionRule{}
	if _, err := f.AbstractComponent.Db.GetLatest(f.Rule); err != nil {
		return err
	}
	f.TopicDeterminer = determiner
	return nil
}

func (f *Actor) Handle(message []byte) ([]byte, error) {
	var (
		m            = map[string]interface{}{}
		messageValue interface{}
		in           bool
	)
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("actor: can't unmarshall message")
	}
	if messageValue, in = m[f.Rule.Key]; !in {
		return nil, errors.New(fmt.Sprintf("actor: no field %s in message %s", f.Rule.Key, string(message)))
	}
	f.Kafka.SetProducerBootstrapServers(strings.Split(f.Rule.BootstrapServers, ","))
	f.Kafka.SetProducerTopics([]string{f.TopicDeterminer(messageValue.(string))})
	return json.Marshal(&m)
}

func (f *Actor) Run() {
	log.Println("actor running")
	go f.Observe()
	f.Kafka.HandleEvents(f.Handle)
}
