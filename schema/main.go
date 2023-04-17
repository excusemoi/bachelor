package main

import (
	"github.com/bachelor/internal/config"
	"github.com/go-pg/pg/v10"
	migrations "github.com/robinjoseph08/go-pg-migrations/v3"
	"github.com/spf13/viper"
	"log"
	"os"
)

func main() {

	var (
		vp  *viper.Viper
		err error
	)

	if vp, err = config.InitConfig("", "config"); err != nil {
		log.Fatal(err)
	}

	conn := pg.Connect(&pg.Options{
		User:     vp.GetString("postgres.login"),
		Password: vp.GetString("postgres.password"),
		Addr:     vp.GetString("postgres.host") + ":" + viper.GetString("postgres.port"),
		Database: vp.GetString("postgres.name"),
	})

	if err = migrations.Run(conn, "./", os.Args); err != nil {
		log.Fatal(err)
	}
}
