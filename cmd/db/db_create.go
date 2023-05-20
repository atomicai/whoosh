package main

import (
	"fmt"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

func main() {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "whoosh",
	})

	if err != nil {
		panic(err)
	}

	fmt.Println(session)
}
