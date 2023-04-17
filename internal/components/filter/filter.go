package filter

import (
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model/filter"
)

type Filter struct {
	fr *filter.FiltrationRule
	components.AbstractComponent
}

func (f *Filter) Init(configPath string, fr *filter.FiltrationRule) error {
	if err := f.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	f.fr = fr
	return nil
}

func (f *Filter) Handle(message []byte) error {
	fmt.Println("filter")
	return nil
}

func (f *Filter) Run() {
	f.Kafka.HandleEvents(f.Handle)
}
