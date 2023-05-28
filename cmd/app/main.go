package main

import (
	"fmt"
	"github.com/atomicai/whoosh/internal/handler"
)

func main() {
	path := handler.InitDijkstra(1, 16)
	fmt.Printf("paths size is: %f", path)
	//handler.OptimalPath()
}
