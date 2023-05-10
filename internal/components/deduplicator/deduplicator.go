package deduplicator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
	"time"
)

type Deduplicator struct {
	components.AbstractComponent[model.DeduplicationRule]
}

func (d *Deduplicator) Init(path string, name string) error {
	return d.AbstractComponent.Init(path, name)
}

func (d *Deduplicator) Handle(message []byte) ([]byte, error) {
	var (
		m   = map[string]interface{}{}
		res []byte
		err error
	)
	for _, rule := range d.Rules {
		if err = json.Unmarshal(message, &m); err != nil {
			return nil, errors.New("deduplicator: can't unmarshall message")
		}

		var (
			currentTime        = time.Now()
			isInsideTimeWindow = currentTime.After(rule.Start) && currentTime.Before(rule.End)
			isKeyExist         = false
			messageValue       interface{}
		)

		if messageValue, isKeyExist = m[rule.Field]; isKeyExist {
			if isInsideTimeWindow {
				if _, in := rule.Value[messageValue.(string)]; in {
					return nil, errors.New(fmt.Sprintf(`deduplicator: message %s didn't pass'`, string(message)))
				}
			}
		}
		if isKeyExist {
			rule.Value = make(map[string]struct{})
			rule.Value[messageValue.(string)] = struct{}{}
		}
		if res, err = json.Marshal(&m); err != nil {
			return nil, errors.New("deduplicator: can't marshal message")
		}
	}
	return res, err
}

func (d *Deduplicator) Run() {
	log.Println("deduplicator running")
	go d.Observe()
	d.Kafka.HandleEvents(d.Handle)
}
