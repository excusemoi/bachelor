package main

import (
	"github.com/bachelor/pkg/config"
	"github.com/bachelor/pkg/db"
	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {
	var (
		err error
	)

	if err = config.InitConfig(""); err != nil {
		log.Fatal(err)
	}

	log.Println("Migrations")

	conn := pg.Connect(&pg.Options{
		User:     viper.GetString("postgres.login"),
		Password: db.GetEnv("POSTGRES_PASSWORD", "root"),
		Addr:     viper.GetString("postgres.host") + ":" + viper.GetString("postgres.port"),
		Database: viper.GetString("postgres.name"),
	})

	if err = migrations.Run(conn, "./", os.Args); err != nil {
		log.Fatal(err)
	}
}
