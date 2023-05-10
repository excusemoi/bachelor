package enricher

import (
	"encoding/json"
	"errors"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
)

type Enricher struct {
	components.AbstractComponent[model.EnrichmentRule]
}

func (r *Enricher) Init(path string, name string) error {
	return r.AbstractComponent.Init(path, name)
}

func (r *Enricher) Handle(message []byte) ([]byte, error) {
	var (
		m   = map[string]interface{}{}
		res []byte
		err error
	)
	for _, rule := range r.Rules {
		if err := json.Unmarshal(message, &m); err != nil {
			return nil, errors.New("enricher: can't unmarshal message")
		}
		m[rule.Field] = rule.Value
		if res, err = json.Marshal(&m); err != nil {
			return nil, errors.New("enricher: can't marshal message")
		}
	}
	return res, err
}

func (r *Enricher) Run() {
	log.Println("enricher running")
	go r.Observe()
	r.Kafka.HandleEvents(r.Handle)
}
