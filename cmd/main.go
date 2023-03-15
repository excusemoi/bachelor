package main

import (
	"github.com/bachelor/pkg/config"
	"github.com/bachelor/pkg/db"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	var (
		dbClient *db.Db
		err      error
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

	log.Println(dbClient)

	for {

	}
}
