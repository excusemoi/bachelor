package transformer

import (
	"fmt"
	"github.com/bachelor/internal/components"
	"github.com/bachelor/internal/model/filter"
)

type Transformer struct {
	components.AbstractComponent
}

func (t *Transformer) Init(configPath string, fr *filter.FiltrationRule) error {
	if err := t.AbstractComponent.Init(configPath); err != nil {
		return err
	}
	//t.fr = fr TODO transformation rule
	return nil
}

func (t *Transformer) Handle(message []byte) error {
	fmt.Println("transformer")
	return nil
}

func (t *Transformer) Run() {
	t.Kafka.HandleEvents(t.Handle)
}
