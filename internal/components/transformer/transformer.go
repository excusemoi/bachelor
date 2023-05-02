package transformer

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
)

type Transformer struct {
	components.AbstractComponent[model.TransformationRule]
}

func (t *Transformer) Init(configPath string) error {
	if err := t.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	t.Rule = &model.TransformationRule{}
	if _, err := t.AbstractComponent.Db.GetLatest(t.Rule); err != nil {
		return err
	}
	return nil
}

func (t *Transformer) Handle(message []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("transformer: can't unmarshall message")
	}
	if v, in := m[t.Rule.Field]; in {
		m[t.Rule.TargetField] = v
		delete(m, t.Rule.Field)
	} else {
		return nil, errors.New(fmt.Sprintf(`transformer: no field %s in message %s`, t.Rule.Field,
			string(message)))
	}
	if res, err := json.Marshal(&m); err != nil {
		return nil, errors.New("transformer: can't marshal message")
	} else {
		return res, nil
	}
}

func (t *Transformer) Run() {
	log.Printf("transformer: running")
	go t.Observe()
	t.Kafka.HandleEvents(t.Handle)
}
