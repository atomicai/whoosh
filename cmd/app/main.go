package main

import (
	"github.com/atomicai/whoosh/internal/handler"
)

func main() {
	//db := database_init.NewDBHandler("whoosh")
	//db.DeleteTables()

	//start := time.Now()
	//db.CreateTables()
	//elapsed := time.Since(start)
	//fmt.Printf("CreateTables function take %s time\n", elapsed)

	//dbname := "whoosh"
	//handler.NewDijkstra(dbname)
	handler.OptimalPath()

	//pathQuery := models.PathQuery{
	//	StartPoint: models.Point{
	//		Lat: 55.69811213203434,
	//		Lon: 37.52609142135623,
	//	},
	//	EndPoint: models.Point{
	//		Lat: 55.69804142135623,
	//		Lon: 37.52626213203435,
	//	},
	//	UserId: "1",
	//}
	//res := handler.Dijkstra(&pathQuery)
	//fmt.Printf("result: %+v", res)
}
