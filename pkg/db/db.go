package db

import (
	"github.com/go-pg/pg/v10"
	"os"
)

type Db struct {
	client *pg.DB
}

func New(dbLogin, dbPass, dbHost, dbPort, dbName string) (*Db, error) {

	connectOptions, err := pg.ParseURL("postgres://" +
		dbLogin + ":" +
		dbPass + "@" +
		dbHost + ":" +
		dbPort + "/" +
		dbName + "?sslmode=disable")
	if err != nil {
		return nil, err
	}

	dbClient := pg.Connect(connectOptions)

	newClient := &Db{
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
