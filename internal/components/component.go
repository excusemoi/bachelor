package components

import (
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/db"
	"github.com/bachelor/internal/kafka"
	"github.com/spf13/viper"
	"sync"
)

type IComponent interface {
	Run()
	GetConsumerTopics() []string
	SetProducerTopics([]string)
	Handle(message []byte) ([]byte, error)
	GetNext() IComponent
	SetNext(next IComponent) IComponent
}

type AbstractComponent[T db.Model] struct {
	Next  IComponent
	Kafka *kafka.Client
	Db    *db.Db[T]
	Rules []*T
	Wg    *sync.WaitGroup
}

func (ac *AbstractComponent[T]) Init(path string, name string) error {
	var (
		vp  *viper.Viper
		err error
	)
	if vp, err = config.InitConfig(path, name); err != nil {
		return err
	}

	ac.Kafka = &kafka.Client{}
	ac.Kafka.Init(vp)

	ac.Rules = []*T{}

	ac.Db = &db.Db[T]{}
	if _, err = ac.Db.GetAll(&ac.Rules); err != nil {
		return err
	}

	if ac.Db, err = ac.Db.Init(vp); err != nil {
		return err
	}

	return nil
}

func (ac *AbstractComponent[T]) SetNext(next IComponent) IComponent {
	if ac.Kafka != nil {
		ac.SetProducerTopics(next.GetConsumerTopics())
	}
	ac.Next = next
	return ac.Next
}

func (ac *AbstractComponent[T]) GetNext() IComponent {
	return ac.Next
}

func (ac *AbstractComponent[T]) RunPipeline() {
	curr := ac.Next
	for curr != nil {
		ac.Wg.Add(1)
		go curr.Run()
		curr = curr.GetNext()
	}
	ac.Wg.Wait()
}

func (ac *AbstractComponent[T]) Observe() {
	ac.Db.Observe((*T)(nil))
}

func (ac *AbstractComponent[T]) SetProducerTopics(topics []string) {
	if ac.Kafka != nil {
		ac.Kafka.SetProducerTopics(topics)
	}
}

func (ac *AbstractComponent[T]) GetConsumerTopics() []string {
	if ac.Kafka != nil {
		return ac.Kafka.GetConsumerTopics()
	}
	return nil
}
