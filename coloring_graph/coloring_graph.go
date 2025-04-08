package coloringgraph

import (
	"github.com/hvuhsg/zkp/graph"
)

type ColoringGraph struct {
	*graph.Graph[ColorNodeValue]
	Colors map[string]struct{}
}

func NewColoringGraph() *ColoringGraph {
	graph := graph.NewGraph[ColorNodeValue]()
	return &ColoringGraph{
		Graph:  graph,
		Colors: make(map[string]struct{}),
	}
}

func (cg *ColoringGraph) Clone() *ColoringGraph {
	return &ColoringGraph{
		Graph:  cg.Graph.Clone(),
		Colors: cg.Colors,
	}
}

func (cg *ColoringGraph) IsGraphColoringValid() bool {
	for _, node := range cg.GetNodes() {
		if _, ok := cg.Colors[string(node.Value)]; !ok {
			return false
		}
	}

	for _, edge := range cg.GetEdges() {
		node1 := cg.GetNodes()[edge.From]
		node2 := cg.GetNodes()[edge.To]

		// Check if nodes have the same color
		if node1.Value == node2.Value {
			return false
		}
	}
	return true
}
