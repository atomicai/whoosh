package main

import "github.com/atomicai/whoosh/internal/database_init"

func main() {
	db := database_init.NewDBHandler("whoosh")
	db.DeleteTables()
	db.CreateTables()
}
