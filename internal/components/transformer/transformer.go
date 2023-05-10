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

func (t *Transformer) Init(path string, name string) error {
	return t.AbstractComponent.Init(path, name)
}

func (t *Transformer) Handle(message []byte) ([]byte, error) {
	var (
		res []byte
		m   = map[string]interface{}{}
		err error
	)
	if err = json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("transformer: can't unmarshall message")
	}
	for _, rule := range t.Rules {
		if v, in := m[rule.Field]; in {
			m[rule.TargetField] = v
			delete(m, rule.Field)
		} else {
			return nil, errors.New(fmt.Sprintf(`transformer: no field %s in message %s`, rule.Field,
				string(message)))
		}
		if res, err = json.Marshal(&m); err != nil {
			return nil, errors.New("transformer: can't marshal message")
		}
	}
	return res, err
}

func (t *Transformer) Run() {
	log.Printf("transformer: running")
	go t.Observe()
	t.Kafka.HandleEvents(t.Handle)
}
