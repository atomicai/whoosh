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
	//		Lat: 0,
	//		Lon: 7.5,
	//	},
	//	EndPoint: models.Point{
	//		Lat: -4.5,
	//		Lon: -5,
	//	},
	//	UserId: "1",
	//}
	//res := handler.Dijkstra(&pathQuery)
	//fmt.Printf("result: %+v", res)
}
