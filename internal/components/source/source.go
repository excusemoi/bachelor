package source

import (
	"encoding/json"
	"errors"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
)

type Source struct {
	components.AbstractComponent[model.AbstractRule]
}

func (s *Source) Init(configPath string) error {
	if err := s.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	return nil
}

func (s *Source) Handle(message []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("source: can't unmarshall message")
	}
	return json.Marshal(&m)
}

func (s *Source) Run() {
	log.Println("source running")
	go s.Observe()
	s.Kafka.HandleEvents(s.Handle)
}
