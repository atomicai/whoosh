package handler

import (
	"github.com/atomicai/whoosh/internal/models"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
	"sync"
)

type dijkstraStruct struct {
	session *r.Session
}

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
		if dist(nodePoint, point) < 0.0001 {
			return node.Id
		}
	}
	return -1
}
