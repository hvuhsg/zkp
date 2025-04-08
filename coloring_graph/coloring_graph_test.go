package coloringgraph

import (
	"testing"

	"github.com/hvuhsg/zkp/graph"
)

func TestValidGraphColoring(t *testing.T) {
	tests := []struct {
		name     string
		graph    *ColoringGraph
		colors   []string
		expected bool
	}{
		{
			name: "empty graph",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				return cg
			}(),
			colors:   []string{},
			expected: true,
		},
		{
			name: "single node",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				return cg
			}(),
			colors:   []string{"red", "blue"},
			expected: true,
		},
		{
			name: "two nodes with different colors",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddEdge(0, 1)
				return cg
			}(),
			colors:   []string{"red", "blue"},
			expected: true,
		},
		{
			name: "triangle with valid coloring",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("green")
				cg.AddEdge(0, 1)
				cg.AddEdge(1, 2)
				cg.AddEdge(2, 0)
				return cg
			}(),
			colors:   []string{"red", "blue", "green"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add colors to the graph's color set
			for _, color := range tt.colors {
				tt.graph.Colors[color] = struct{}{}
			}

			// Test the coloring validation
			result := tt.graph.IsGraphColoringValid()
			if result != tt.expected {
				t.Errorf("IsGraphColoringValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInvalidGraphColoring(t *testing.T) {
	tests := []struct {
		name     string
		graph    *ColoringGraph
		colors   []string
		expected bool
	}{
		{
			name: "two nodes with same color",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("red")
				cg.AddEdge(0, 1)
				return cg
			}(),
			colors:   []string{"red", "blue"},
			expected: false,
		},
		{
			name: "triangle with invalid coloring",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("red") // Same color as node 0
				cg.AddEdge(0, 1)
				cg.AddEdge(1, 2)
				cg.AddEdge(2, 0)
				return cg
			}(),
			colors:   []string{"red", "blue"},
			expected: false,
		},
		{
			name: "node with color not in allowed set",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red") // Color 3 not in allowed set
				return cg
			}(),
			colors:   []string{"green", "blue"},
			expected: false,
		},
		{
			name: "complex graph with invalid coloring",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				// Create a 4-node graph with a cycle
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("red") // Same color as node 0
				cg.AddNode("green")
				cg.AddEdge(0, 1)
				cg.AddEdge(1, 2)
				cg.AddEdge(2, 3)
				cg.AddEdge(3, 0)
				cg.AddEdge(2, 0) // Same color nodes
				return cg
			}(),
			colors:   []string{"red", "blue", "green"},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add colors to the graph's color set
			for _, color := range tt.colors {
				tt.graph.Colors[color] = struct{}{}
			}

			// Test the coloring validation
			result := tt.graph.IsGraphColoringValid()
			if result != tt.expected {
				t.Errorf("IsGraphColoringValid() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestColoringGraphCreation(t *testing.T) {
	tests := []struct {
		name     string
		graph    *ColoringGraph
		expected bool
	}{
		{
			name:     "new empty graph",
			graph:    NewColoringGraph(),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.graph == nil {
				t.Error("NewColoringGraph() returned nil")
			}
			if tt.graph.Graph == nil {
				t.Error("NewColoringGraph() created graph with nil Graph field")
			}
			if tt.graph.Colors == nil {
				t.Error("NewColoringGraph() created graph with nil Colors map")
			}
		})
	}
}

func TestColoringGraphSerialization(t *testing.T) {
	tests := []struct {
		name     string
		graph    *ColoringGraph
		expected []byte
	}{
		{
			name:     "empty graph",
			graph:    NewColoringGraph(),
			expected: []byte{},
		},
		{
			name: "graph with one node",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				return cg
			}(),
		},
		{
			name: "graph with two nodes",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddEdge(0, 1)
				return cg
			}(),
		},
		{
			name: "graph with three nodes",
			graph: func() *ColoringGraph {
				cg := NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("green")
				cg.AddEdge(0, 1)
				cg.AddEdge(1, 2)
				cg.AddEdge(2, 0)
				return cg
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := tt.graph.Serialize()
			deserialized, err := graph.DeserializeGraph(serialized, DeserializeColorNodeValue)
			if err != nil {
				t.Errorf("error deserializing graph: %v", err)
			}
			// iterate over the nodes and edges of the deserialized graph and compare them to the original graph
			for i, node := range tt.graph.Graph.GetNodes() {
				otherNode := deserialized.GetNodes()[i]
				if err != nil {
					t.Errorf("error getting node %d: %v", i, err)
				}
				if node.Id != otherNode.Id {
					t.Errorf("node ids %d does not match", i)
				}
				if node.Value != otherNode.Value {
					t.Errorf("node values %d does not match", i)
				}
			}
			for i, edge := range tt.graph.Graph.GetEdges() {
				otherEdge := deserialized.GetEdges()[i]
				if err != nil {
					t.Errorf("error getting edge %d: %v", i, err)
				}
				if edge.From != otherEdge.From {
					t.Errorf("edge from %d does not match", i)
				}
				if edge.To != otherEdge.To {
					t.Errorf("edge to %d does not match", i)
				}
			}
		})
	}
}
