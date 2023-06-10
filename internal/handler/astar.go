package handler

import (
	"fmt"
	"github.com/atomicai/whoosh/internal/models"
	"log"
	"math"
)

func AStar(pathQuery *models.PathQuery) *models.PathResponse {
	graph := NewGraph()
	err := graph.InitGraph(dijkstra.session)
	if err != nil {
		log.Fatal(err)
	}

	startId := graph.GetNodeId(pathQuery.StartPoint)
	finishId := graph.GetNodeId(pathQuery.EndPoint)
	if startId == -1 && finishId == -1 {
		log.Fatalln("no point with this coordinates")
	}

	graph.AStarHelper(startId, finishId)
	start := graph.GetNode(startId)
	finish := graph.GetNode(finishId)

	fmt.Printf("dist in astar: %f\n", finish.Value)

	path := make([]models.Point, 0, 10)
	now := finish
	for now != start {
		point := models.Point{Lat: now.Lat, Lon: now.Lon}
		path = append(path, point)
		now = now.Parent
	}
	point := models.Point{Lat: now.Lat, Lon: now.Lon}
	path = append(path, point)

	result := models.PathResponse{Path: make([]models.Point, 0, len(path)), UserId: pathQuery.UserId}

	for i := len(path) - 1; i >= 0; i-- {
		result.Path = append(result.Path, path[i])
	}

	return &result
}

func (g *Graph) AStarHelper(startId, finishId int) {
	visited := make(map[int]bool)
	HevHeap := &models.HevHeap{}
	finishNode := g.GetNode(finishId)

	startNode := g.GetNode(startId)
	startNode.Value = 0
	startNode.HValue = startNode.Value + hevristic(startNode, finishNode)
	HevHeap.Push(startNode)

	for HevHeap.Size() > 0 {
		current := HevHeap.Pop()
		if current.Id == finishId {
			break
		}

		visited[current.Id] = true
		edges := current.Neighbors

		for _, edge := range edges {
			to := g.GetNode(edge.To)
			if current.Value+edge.Weight < to.Value {
				to.Value = current.Value + edge.Weight
				to.HValue = to.Value + hevristic(to, finishNode)
				to.Parent = current
				HevHeap.Push(to)
			}
		}
	}
}

func hevristic(start, end *models.NodeWithValue) float64 {
	return (math.Abs(start.Lat-end.Lat) + math.Abs(start.Lon-end.Lon)) * 1000 // check param 1000!!!
	//return math.Sqrt((start.Lat-end.Lat)*(start.Lat-end.Lat) + (start.Lon-end.Lon)*(start.Lon-end.Lon))
}
