package main

import (
	"github.com/atomicai/whoosh/internal/handler"
)

func main() {
	dbname := "whoosh"
	handler.NewDijkstra(dbname)
	handler.OptimalPath()

	//pathQuery := models.PathQuery{
	//	StartPoint: models.Point{
	//		Lat: 37.52626213203435,
	//		Lon: 55.69804142135623,
	//	},
	//	EndPoint: models.Point{
	//		Lat: 37.526432842712474,
	//		Lon: 55.69811213203434,
	//	},
	//	UserId: "1",
	//}
	//res := handler.Dijkstra(&pathQuery)
	//fmt.Printf("result: %+v", res)
}
