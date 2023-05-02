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

func (d *Deduplicator) Init(configPath string) error {
	if err := d.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	d.Rule = &model.DeduplicationRule{}
	if _, err := d.AbstractComponent.Db.GetLatest(d.Rule); err != nil {
		return err
	}
	return nil
}

func (d *Deduplicator) Handle(message []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("deduplicator: can't unmarshall message")
	}

	var (
		currentTime        = time.Now()
		isInsideTimeWindow = currentTime.After(d.Rule.Start) && currentTime.Before(d.Rule.End)
		isKeyExist         = false
		messageValue       interface{}
	)

	if messageValue, isKeyExist = m[d.Rule.Field]; isKeyExist {
		if isInsideTimeWindow {
			if _, in := d.Rule.Value[messageValue.(string)]; in {
				return nil, errors.New(fmt.Sprintf(`deduplicator: message %s didn't pass'`, string(message)))
			}
		}
	}
	if isKeyExist {
		d.Rule.Value = make(map[string]struct{})
		d.Rule.Value[messageValue.(string)] = struct{}{}
	}
	if res, err := json.Marshal(&m); err != nil {
		return nil, errors.New("deduplicator: can't marshal message")
	} else {
		return res, nil
	}
}

func (d *Deduplicator) Run() {
	log.Println("deduplicator running")
	go d.Observe()
	d.Kafka.HandleEvents(d.Handle)
}
