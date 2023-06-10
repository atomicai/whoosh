package main

import (
	"fmt"
	"github.com/atomicai/whoosh/internal/database_init"
	"github.com/atomicai/whoosh/internal/handler"
	"time"
)

func main() {
	db := database_init.NewDBHandler("test")
	db.DeleteTables()

	start := time.Now()
	db.CreateTables()
	elapsed := time.Since(start)
	fmt.Printf("CreateTables function take %s time\n", elapsed)

	dbname := "test"
	handler.NewDijkstra(dbname)
	handler.OptimalPath()

	//pathQuery := models.PathQuery{
	//	StartPoint: models.Point{
	//		Lat: 55.70094055915914,
	//		Lon: 37.54008969696205,
	//	},
	//	EndPoint: models.Point{
	//		Lat: 55.70101126983727,
	//		Lon: 37.52762781745931,
	//	},
	//	UserId: "1",
	//}
	//res := handler.Dijkstra(&pathQuery)
	//fmt.Printf("result: %+v", res)
}
