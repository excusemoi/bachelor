package main

import (
	"context"
	"github.com/bachelor/pkg/config"
	"github.com/bachelor/pkg/db"
	"github.com/bachelor/pkg/model"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"sync"
	"time"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		dbClient    *db.Db
		cfr         = &model.FiltrationRule{Mx: sync.RWMutex{}}
		err         error
	)

	if err = config.InitConfig(""); err != nil {
		log.Fatal(err)
	}

	if dbClient, err = db.New(
		viper.GetString("postgres.login"),
		os.Getenv("POSTGRES_PASSWORD"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.name"),
	); err != nil {
		log.Fatal(err)
	}

	go handleFiltrationRuleChanges(ctx, dbClient, cfr)

	//just for loop
	for {
	}

	cancel()
}

func handleFiltrationRuleChanges(ctx context.Context, dbClient *db.Db, cfr *model.FiltrationRule) {
	for {
		select {
		case <-time.After(time.Second * 10):
			{
				nfr, _ := dbClient.GetLatestFiltrationRule()
				if !reflect.DeepEqual(nfr, cfr) {
					cfr.Mx.Lock()
					cfr = nfr
					cfr.Mx.Unlock()
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
