package main

import (
	"fmt"
	"github.com/atomicai/whoosh/internal/database_init"
	"time"
)

func main() {
	db := database_init.NewDBHandler("whoosh")
	db.DeleteTables()

	start := time.Now()
	db.CreateTables()
	elapsed := time.Since(start)
	fmt.Printf("CreateTables function take %s time", elapsed)
}
