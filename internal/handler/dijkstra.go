package handler

import (
	"encoding/csv"
	"fmt"
	"github.com/atomicai/whoosh/internal/models"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"sync"
)

type Graph struct {
	Nodes []*models.Node
	Edges map[*models.Node][]*models.Edge
	Mp    map[int]*models.Node
	sync.Mutex
}

func NewGraph() *Graph {
	return &Graph{Nodes: make([]*models.Node, 0), Edges: make(map[*models.Node][]*models.Edge), Mp: make(map[int]*models.Node)}
}

func (g *Graph) GetNode(id int) (node *models.Node) {
	g.Lock()
	defer g.Unlock()
	return g.Mp[id]
}

func (g *Graph) AddNode(node models.Node) {
	g.Lock()
	defer g.Unlock()
	g.Nodes = append(g.Nodes, &node)
	g.Mp[node.Id] = &node
}

func (g *Graph) AddEdge(n1, n2 *models.Node, weight float64) {
	g.Lock()
	defer g.Unlock()
	g.Edges[n1] = append(g.Edges[n1], &models.Edge{Node: n2, Weight: weight})
	g.Edges[n2] = append(g.Edges[n2], &models.Edge{Node: n1, Weight: weight})
}

func InitDijkstra() {
	graph1 := NewGraph()
	graph1.InitNodes()
	graph1.InitEdges()
	graph1.Dijkstra(5)
	fmt.Println(graph1.Mp[7].Value)

	graph2 := NewGraph()
	graph2.InitNodes()
	graph2.InitEdges()
	graph2.AStar(5, 7)
	fmt.Println(graph2.Mp[7].Value)
}

func (g *Graph) InitNodes() {
	dirLink := "C:\\Users\\insha\\OneDrive\\Документы\\whoosh\\TestDatasets\\nodes.csv" // path in computer
	csvFile, err := os.Open(dirLink)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	i := 0
	reader := csv.NewReader(csvFile)
	for ; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if i == 0 {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		node, err := parseNode(line)
		if err != nil {
			log.Fatal(err)
		}

		g.AddNode(node.(models.Node))
	}
}

func parseNode(line []string) (interface{}, error) {
	id, err := strconv.Atoi(line[0])
	if err != nil {
		return nil, err
	}

	x, err := strconv.ParseFloat(line[1], 64)
	if err != nil {
		return nil, err
	}

	y, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return nil, err
	}
	node := models.Node{Id: id, X: x, Y: y, Value: 1e9}
	return node, nil
}

func (g *Graph) InitEdges() {
	dirLink := "C:\\Users\\insha\\OneDrive\\Документы\\whoosh\\TestDatasets\\edges.csv" // path in computer
	csvFile, err := os.Open(dirLink)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	i := 0
	reader := csv.NewReader(csvFile)
	for ; ; i++ {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if i == 0 {
			continue
		}
		if err != nil {
			log.Fatal(err)
		}

		err = g.parseEdge(line)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (g *Graph) parseEdge(line []string) error {
	from, err := strconv.Atoi(line[0])
	if err != nil {
		return err
	}

	to, err := strconv.Atoi(line[1])
	if err != nil {
		return err
	}

	weight, err := strconv.ParseFloat(line[2], 64)
	if err != nil {
		return err
	}
	g.AddEdge(g.Mp[from], g.Mp[to], weight)
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
		edges := g.Edges[current]

		for _, edge := range edges {
			if !visited[edge.Node.Id] {
				if current.Value+edge.Weight < edge.Node.Value {
					edge.Node.Value = current.Value + edge.Weight
					edge.Node.Parent = current
					ValueHeap.Push(edge.Node)
				}
			}
		}
	}
}

func (g *Graph) AStar(startId, finishId int) {
	visited := make(map[int]bool)
	HevHeap := &models.HevHeap{}

	startNode := g.GetNode(startId)
	startNode.Value = 0
	startNode.HevristicValue = startNode.Value + hevristic(g.Mp[startId], g.Mp[finishId])
	HevHeap.Push(startNode)

	for HevHeap.Size() > 0 {
		current := HevHeap.Pop()
		if current.Id == finishId {
			break
		}

		visited[current.Id] = true
		edges := g.Edges[current]

		for _, edge := range edges {
			if current.Value+edge.Weight < edge.Node.Value {
				edge.Node.Value = current.Value + edge.Weight
				edge.Node.HevristicValue = edge.Node.Value + hevristic(edge.Node, g.Mp[finishId])
				edge.Node.Parent = current
				HevHeap.Push(edge.Node)
			}
		}
	}
}

func hevristic(start, end *models.Node) float64 {
	return (math.Abs(start.X-end.X) + math.Abs(start.Y-end.Y)) / 100 // check param 100!!!
}
