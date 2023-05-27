package models

type Node struct {
	Id             int
	X              float64
	Y              float64
	Value          float64
	HevristicValue float64
	Parent         *Node
}

type Edge struct {
	Node   *Node
	Weight float64
}

type Coordinates struct {
	X int
	Y int
}
