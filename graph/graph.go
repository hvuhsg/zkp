package graph

const (
	version = 0b00000001
)

type NodeValue interface {
	Serialize() []byte
}

type Graph[T NodeValue] struct {
	nodes []Node[T]
	edges []Edge
}

type Node[T NodeValue] struct {
	Id    uint16
	Value T
}

type Edge struct {
	From  int
	To    int
	Nonce int
}

func NewGraph[T NodeValue]() *Graph[T] {
	return &Graph[T]{
		nodes: make([]Node[T], 0),
		edges: make([]Edge, 0),
	}
}

func (g *Graph[T]) AddNode(value T) {
	id := uint16(len(g.nodes))
	g.nodes = append(g.nodes, Node[T]{Id: id, Value: value})
}

func (g *Graph[T]) AddEdge(from, to int) {
	g.edges = append(g.edges, Edge{From: from, To: to})
}

func (g *Graph[T]) GetNodes() []Node[T] {
	return g.nodes
}

func (g *Graph[T]) GetEdges() []Edge {
	return g.edges
}
