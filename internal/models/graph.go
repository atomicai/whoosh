package models

type Edge struct {
	To     int
	Weight float64
}

type InputNode struct {
	Id        int
	Lat       float64
	Lon       float64
	Neighbors []*Edge
	id        string
}

type Node struct {
	NodeId    int
	Lat       float64
	Lon       float64
	Neighbors []*Edge
}

type NodeWithValue struct {
	Id        int
	Lat       float64
	Lon       float64
	Neighbors []*Edge
	Value     float64
	HValue    float64
	Parent    *NodeWithValue
}

func NewNode() *Node {
	return &Node{Neighbors: make([]*Edge, 0)}
}

func (n *Node) AddNeighbor(edge *Edge) {
	n.Neighbors = append(n.Neighbors, edge)
}
