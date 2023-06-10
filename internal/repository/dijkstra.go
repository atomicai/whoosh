package repository

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
)

type IDijkstraRepository interface {
}

type DijkstraRepository struct {
	dbName  string
	session *r.Session
}

func NewDijkstraRepository(dbname string) *DijkstraRepository {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "rethinkdb",
		Database: dbname,
	})

	if err != nil {
		fmt.Println("error on connecting to database")
		log.Fatal(err)
	}

	return &DijkstraRepository{dbName: dbname, session: session}
}
