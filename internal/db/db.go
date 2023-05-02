package db

import (
	"github.com/bachelor/internal/model"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/viper"
	"os"
	"sync"
)

type Model interface {
	model.FiltrationRule |
		model.TransformationRule |
		model.DeduplicationRule |
		model.EnrichmentRule |
		model.ActionRule |
		model.AbstractRule
}

type Db[T Model] struct {
	client *pg.DB
	mx     sync.Mutex
}

func (db *Db[T]) Init(vp *viper.Viper) (*Db[T], error) {

	connectOptions, err := pg.ParseURL("postgres://" +
		vp.GetString("postgres.login") + ":" +
		GetEnv("POSTGRES_PASSWORD", "root") + "@" +
		vp.GetString("postgres.host") + ":" +
		vp.GetString("postgres.port") + "/" +
		vp.GetString("postgres.name") + "?sslmode=disable")
	if err != nil {
		return nil, err
	}

	dbClient := pg.Connect(connectOptions)

	newClient := &Db[T]{
		client: dbClient,
	}

	return newClient, nil
}

func GetEnv(key, def string) string {
	e := os.Getenv(key)
	if e == "" {
		return def
	}
	return e
}
