package filter

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model"
	"log"
	"strings"
)

const (
	equals      = "equals"
	notEquals   = "notEquals"
	contains    = "contains"
	notContains = "notContains"
)

type Filter struct {
	components.AbstractComponent[model.FiltrationRule]
}

func (f *Filter) Init(configPath string) error {
	if err := f.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	f.Rule = &model.FiltrationRule{}
	if _, err := f.AbstractComponent.Db.GetLatest(f.Rule); err != nil {
		return err
	}
	return nil
}

func (f *Filter) Handle(message []byte) ([]byte, error) {
	m := map[string]interface{}{}
	if err := json.Unmarshal(message, &m); err != nil {
		return nil, errors.New("filter: can't unmarshall message")
	}
	if value, in := m[f.Rule.Field]; in {
		if ok, err := f.FilterFunction(value.(string)); err != nil {
			return nil, err
		} else if !ok {
			return nil, errors.New(fmt.Sprintf("filter: message %s didn't pass", string(message)))
		}
	} else {
		return nil, errors.New(fmt.Sprintf(`filter: no field %s in message %s`, f.Rule.Field,
			string(message)))
	}
	if res, err := json.Marshal(&m); err != nil {
		return nil, errors.New("filter: can't marshal message")
	} else {
		return res, nil
	}
}

func (f *Filter) Run() {
	log.Println("filter running")
	go f.Observe()
	f.Kafka.HandleEvents(f.Handle)
}
func (f *Filter) FilterFunction(value string) (bool, error) {
	switch f.AbstractComponent.Rule.Function {
	case equals:
		return f.equals(value), nil
	case notEquals:
		return !f.equals(value), nil
	case contains:
		return f.contains(value), nil
	case notContains:
		return !f.contains(value), nil
	}
	return false, errors.New("incorrect function name")
}

func (f *Filter) equals(value string) bool {
	return f.Rule.Value == value
}

func (f *Filter) contains(value string) bool {
	return strings.Contains(value, f.AbstractComponent.Rule.Value)
}
