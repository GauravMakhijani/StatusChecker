package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	uri      = "postgres://gaurav:4444@localhost:5432/websites?sslmode=disable"
)

var DB *sqlx.DB

func Init() {

	var err error

	DB, err = sqlx.Connect(dbDriver, uri)
	if err != nil {
		log.Fatal("Cannot initialize database", err)
	}

	log.Println("Database initialized")
}
