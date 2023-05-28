package handler

import (
	"github.com/atomicai/whoosh/internal/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"log"
	"math"
	"sync"
)

type Graph struct {
	Nodes []*models.NodeWithValue
	Mp    map[int]*models.NodeWithValue
	sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{Nodes: make([]*models.NodeWithValue, 0), Mp: make(map[int]*models.NodeWithValue)}
}

func (g *Graph) GetNode(id int) (node *models.NodeWithValue) {
	g.Lock()
	defer g.Unlock()
	return g.Mp[id]
}

func (g *Graph) AddNode(node models.NodeWithValue) {
	g.Lock()
	defer g.Unlock()
	g.Nodes = append(g.Nodes, &node)
	g.Mp[node.Id] = &node
}

func (g *Graph) GetNodeId(point models.Point) int {
	for _, node := range g.Nodes {
		nodePoint := models.Point{Lat: node.Lat, Lon: node.Lon}
		if dist(nodePoint, point) < 10 {
			return node.Id
		}
	}
	return -1
}

func dist(a, b models.Point) float64 {
	return math.Sqrt((a.Lat-b.Lat)*(a.Lat-b.Lat) + (a.Lon-b.Lon)*(a.Lon-b.Lon))
}

func InitDijkstra(pathQuery *models.PathQuery) *models.PathResponse {
	session, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "whoosh",
	})
	if err != nil {
		log.Fatal(err)
	}

	graph := NewGraph()
	err = graph.InitGraph(session)
	if err != nil {
		log.Fatal(err)
	}

	startId := graph.GetNodeId(pathQuery.StartPoint)
	finishId := graph.GetNodeId(pathQuery.EndPoint)
	if startId == -1 && finishId == -1 {
		log.Fatalln("no point with this coordinates")
	}

	graph.Dijkstra(startId)
	//graph.AStar(startId, finishId)
	start := graph.GetNode(startId)
	finish := graph.GetNode(finishId)

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

func (g *Graph) Dijkstra(Id int) {
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

func (g *Graph) AStar(startId, finishId int) {
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
	return (math.Abs(start.Lat-end.Lat) + math.Abs(start.Lon-end.Lon)) / 100 // check param 100!!!
}
