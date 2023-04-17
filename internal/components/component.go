package components

import (
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/kafka"
	"github.com/spf13/viper"
)

type IComponent interface {
	Run()
	GetConsumerTopics() []string
	SetProducerTopics([]string)
	Handle(message []byte) error
	GetNext() IComponent
}

type AbstractComponent struct {
	IComponent
	Next  IComponent
	Kafka *kafka.Client
}

func (ac *AbstractComponent) Init(configPath string) error {
	var (
		vp  *viper.Viper
		err error
	)
	if vp, err = config.InitConfig(configPath, "config"); err != nil {
		return err
	}
	ac.Kafka = &kafka.Client{}
	ac.Kafka.Init(vp)
	return nil
}

func (ac *AbstractComponent) SetNext(next IComponent) {
	if ac.Kafka != nil {
		ac.SetProducerTopics(next.GetConsumerTopics())
	}
	ac.Next = next
}

func (ac *AbstractComponent) GetNext() IComponent {
	return ac.Next
}

func (ac *AbstractComponent) RunPipeline() {
	curr := ac.Next
	for curr != nil {
		curr.Run()
		curr = curr.GetNext()
	}
}

func (ac *AbstractComponent) SetProducerTopics(topics []string) {
	if ac.Kafka != nil {
		ac.Kafka.SetProducerTopics(topics)
	}
}

func (ac *AbstractComponent) GetConsumerTopics() []string {
	if ac.Kafka != nil {
		return ac.Kafka.GetConsumerTopics()
	}
	return []string{}
}
