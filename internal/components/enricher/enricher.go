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

func (r *Enricher) Init(configPath string) error {
	if err := r.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	r.Rule = &model.EnrichmentRule{}
	if _, err := r.AbstractComponent.Db.GetLatest(r.Rule); err != nil {
		return err
	}
	return nil
}

func (r *Enricher) Handle(message []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("enricher: can't unmarshal message")
	}
	m[r.Rule.Field] = r.Rule.Value
	if res, err := json.Marshal(&m); err != nil {
		return nil, errors.New("enricher: can't marshal message")
	} else {
		return res, nil
	}
}

func (r *Enricher) Run() {
	log.Println("enricher running")
	go r.Observe()
	r.Kafka.HandleEvents(r.Handle)
}
