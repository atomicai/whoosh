package main

import (
	"fmt"
	"github.com/atomicai/whoosh/internal/handler"
	"github.com/atomicai/whoosh/internal/models"
)

func main() {
	pathQuery := models.PathQuery{
		StartPoint: models.Point{
			Lat: -13.5,
			Lon: 0,
		},
		EndPoint: models.Point{
			Lat: 13.5,
			Lon: 0,
		},
		UserId: "4",
	}

	path := handler.InitDijkstra(&pathQuery)
	fmt.Printf("path: %+v", path)
}
