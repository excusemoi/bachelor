package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
)

type Filter struct {
	components.AbstractComponent[model.FiltrationRule]
}

func (f *Filter) Init(path, name string) error {
	if err := f.AbstractComponent.Init(path, name); err != nil {
		return err
	}
	return nil
}

func (f *Filter) Handle(message []byte) ([]byte, error) {
	var (
		m   = map[string]interface{}{}
		res []byte
		err error
	)
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("filter: can't unmarshall message")
	}
	for _, rule := range f.Rules {
		if value, in := m[rule.Field]; in {
			if ok, err := rule.FilterFunction(value.(string)); err != nil {
				return nil, err
			} else if !ok {
				return nil, errors.New(fmt.Sprintf("filter: message %s didn't pass", string(message)))
			}
		} else {
			return nil, errors.New(fmt.Sprintf(`filter: no field %s in message %s`, rule.Field,
				string(message)))
		}
		if res, err = json.Marshal(&m); err != nil {
			return nil, errors.New("filter: can't marshal message")
		}
	}
	return res, err
}

func (f *Filter) Run() {
	log.Println("filter running")
	go f.Observe()
	f.Kafka.HandleEvents(f.Handle)
}
