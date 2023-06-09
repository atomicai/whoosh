package handler

import (
	"fmt"
	"github.com/atomicai/whoosh/internal/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
	"math"
)

var dijkstra dijkstraStruct

func NewDijkstra(dbname string) {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: dbname,
	})
	if err != nil {
		log.Fatal(err)
	}

	dijkstra = dijkstraStruct{
		session: session,
	}
}

func (g *Graph) InitGraph(session *r.Session) error {
	var result []models.Node

	rows, err := r.Table("graph").Run(session)
	if err != nil {
		return err
	}

	err = rows.All(&result)
	if err != nil {
		return err
	}

	for _, nod := range result {
		nodeWithValue := models.NodeWithValue{
			Id:        nod.NodeId,
			Lat:       nod.Lat,
			Lon:       nod.Lon,
			Neighbors: nod.Neighbors,
			Value:     1e9,
			HValue:    1e9,
		}
		g.AddNode(nodeWithValue)
	}
	return nil
}

func dist(a, b models.Point) float64 {
	return math.Sqrt((a.Lat-b.Lat)*(a.Lat-b.Lat) + (a.Lon-b.Lon)*(a.Lon-b.Lon))
}

func Dijkstra(pathQuery *models.PathQuery) *models.PathResponse {
	graph := NewGraph()
	err := graph.InitGraph(dijkstra.session)
	if err != nil {
		log.Fatal(err)
	}

	startId := graph.GetNodeId(pathQuery.StartPoint)
	finishId := graph.GetNodeId(pathQuery.EndPoint)
	if startId == -1 || finishId == -1 {
		fmt.Println("no point with this coordinates")
		return nil
	}

	graph.DijkstraHelper(startId)
	start := graph.GetNode(startId)
	finish := graph.GetNode(finishId)
	fmt.Printf("dist in dijkstra: %f\n", finish.Value)

	path := make([]models.Point, 0)
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

func (g *Graph) DijkstraHelper(Id int) {
	visited := make(map[int]bool)
	ValueHeap := &models.ValueHeap{}

	startNode := g.GetNode(Id)
	startNode.Value = 0
	ValueHeap.Push(startNode)

	for ValueHeap.Size() > 0 {
		current := ValueHeap.Pop()
		if visited[current.Id] {
			continue
		}

		visited[current.Id] = true
		edges := current.Neighbors

		for _, edge := range edges {
			to := g.GetNode(edge.To)
			if !visited[edge.To] {
				if current.Value+edge.Weight < to.Value {
					to.Value = current.Value + edge.Weight
					to.Parent = current
					ValueHeap.Push(to)
				}
			}
		}
	}
}
