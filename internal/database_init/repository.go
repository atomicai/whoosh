package database_init

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
	"sync"
)

type IDBRepository interface {
	DeleteTables()
	CreateTable(tableName string)
	AddRows(values *[]interface{}, tableName string)
	AddRowsByChan(ch chan interface{}, tableName string)
}

type DBRepository struct {
	dbName  string
	session *r.Session
	sync.Mutex
}

func NewDBRepository(dbname string) *DBRepository {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbname,
	})

	if err != nil {
		fmt.Println("error on connecting to database")
		log.Fatal(err)
	}

	return &DBRepository{dbName: dbname, session: session}
}

func (d *DBRepository) DeleteTables() {
	result, err := r.DB("whoosh").TableList().Run(d.session)
	if err != nil {
		log.Fatal(err)
	}

	var response string
	for result.Next(&response) {
		err := r.DB(d.dbName).TableDrop(response).Exec(d.session)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *DBRepository) CreateTable(tableName string) {
	err := r.DB(d.dbName).TableCreate(tableName).Exec(d.session)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DBRepository) AddRows(values *[]interface{}, tableName string) {
	err := r.Table(tableName).Insert(*values).Exec(d.session)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DBRepository) AddRowsByChan(ch chan interface{}, tableName string) {
	d.Lock()
	for value := range ch {
		err := r.Table(tableName).Insert(value).Exec(d.session)
		if err != nil {
			log.Fatal(err)
		}
	}
	d.Unlock()
}
