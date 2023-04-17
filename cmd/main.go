package main

import (
	"context"
	"fmt"
	"github.com/bachelor/internal/config"
	"github.com/bachelor/internal/db"
	"github.com/bachelor/internal/model/filter"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		dbClient    *db.Db
		//kafkaClient *kafka.Client
		vp  *viper.Viper
		cfr = &filter.FiltrationRule{Mx: sync.RWMutex{}}
		err error
	)

	if err = godotenv.Load(filepath.Join(".env")); err != nil {
		log.Fatal(err)
	}

	if vp, err = config.InitConfig("configs", "config"); err != nil {
		log.Fatal(err)
	}

	if dbClient, err = db.New(
		vp.GetString("postgres.login"),
		os.Getenv("POSTGRES_PASSWORD"),
		vp.GetString("postgres.host"),
		vp.GetString("postgres.port"),
		vp.GetString("postgres.name"),
	); err != nil {
		log.Fatal(err)
	}

	go handleFiltrationRuleChanges(ctx, dbClient, cfr) //TODO observer must be defined in every component

	//if kafkaClient, err = kafka.New(viper.GetString("kafka.bs"), viper.GetStringSlice("kafka.topics")); err != nil {
	//	log.Fatal(err)
	//}
	//
	//kafkaClient.Run()

	fmt.Println("Server started")
	for {
	}

	cancel()
}

func handleFiltrationRuleChanges(ctx context.Context, dbClient *db.Db, cfr *filter.FiltrationRule) {
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
