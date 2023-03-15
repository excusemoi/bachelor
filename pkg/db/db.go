package db

import (
	"github.com/go-pg/pg/v10"
)

type Db struct {
	dbClient *pg.DB
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
		dbClient: dbClient,
	}

	return newClient, nil
}
