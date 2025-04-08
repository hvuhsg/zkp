package commitmentgraph

import (
	"strings"
	"testing"

	coloringgraph "github.com/hvuhsg/zkp/coloring_graph"
)

func TestRandomStringGeneration(t *testing.T) {
	tests := []struct {
		name     string
		length   int
		expected int
	}{
		{
			name:     "zero length",
			length:   0,
			expected: 0,
		},
		{
			name:     "positive length",
			length:   10,
			expected: 10,
		},
		{
			name:     "longer length",
			length:   20,
			expected: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generateRandomString(tt.length)
			if len(result) != tt.expected {
				t.Errorf("generateRandomString(%d) length = %d, want %d", tt.length, len(result), tt.expected)
			}
			// Check that all characters are from letterBytes
			for _, c := range result {
				if !strings.Contains(letterBytes, string(c)) {
					t.Errorf("generateRandomString(%d) contains invalid character %c", tt.length, c)
				}
			}
		})
	}
}

func TestNewCommitmentGraph(t *testing.T) {
	tests := []struct {
		name     string
		graph    *coloringgraph.ColoringGraph
		expected int
	}{
		{
			name: "empty graph",
			graph: func() *coloringgraph.ColoringGraph {
				cg := coloringgraph.NewColoringGraph()
				return cg
			}(),
			expected: 0,
		},
		{
			name: "single node",
			graph: func() *coloringgraph.ColoringGraph {
				cg := coloringgraph.NewColoringGraph()
				cg.AddNode("red")
				return cg
			}(),
			expected: 1,
		},
		{
			name: "multiple nodes",
			graph: func() *coloringgraph.ColoringGraph {
				cg := coloringgraph.NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("green")
				return cg
			}(),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := NewCommitmentGraph(tt.graph)
			if cg == nil {
				t.Error("NewCommitmentGraph returned nil")
			}
			if cg.nodesValues == nil {
				t.Fatal("nodesValues is nil")
			}
			if len(cg.nodesValues) != tt.expected {
				t.Errorf("nodesValues length = %d, want %d", len(cg.nodesValues), tt.expected)
			}
			if len(cg.GetNodes()) != tt.expected {
				t.Errorf("GetNodes length = %d, want %d", len(cg.GetNodes()), tt.expected)
			}
		})
	}
}

func TestGetNodeValue(t *testing.T) {
	tests := []struct {
		name     string
		graph    *coloringgraph.ColoringGraph
		nodeID   int
		expected string
	}{
		{
			name: "first node",
			graph: func() *coloringgraph.ColoringGraph {
				cg := coloringgraph.NewColoringGraph()
				cg.AddNode("red")
				return cg
			}(),
			nodeID: 0,
		},
		{
			name: "middle node",
			graph: func() *coloringgraph.ColoringGraph {
				cg := coloringgraph.NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("green")
				return cg
			}(),
			nodeID: 1,
		},
		{
			name: "last node",
			graph: func() *coloringgraph.ColoringGraph {
				cg := coloringgraph.NewColoringGraph()
				cg.AddNode("red")
				cg.AddNode("blue")
				cg.AddNode("green")
				return cg
			}(),
			nodeID: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := NewCommitmentGraph(tt.graph)
			value := cg.GetNodeValue(tt.nodeID)

			// Check that the value contains the original color and a random string
			parts := strings.Split(value, "|")
			if len(parts) != 2 {
				t.Errorf("GetNodeValue(%d) returned invalid format: %s", tt.nodeID, value)
			}

			// Check that the first part matches the original color
			originalColor := string(tt.graph.GetNodes()[tt.nodeID].Value)
			if parts[0] != originalColor {
				t.Errorf("GetNodeValue(%d) color = %s, want %s", tt.nodeID, parts[0], originalColor)
			}

			// Check that the second part is a valid random string
			if len(parts[1]) != 10 {
				t.Errorf("GetNodeValue(%d) random string length = %d, want 10", tt.nodeID, len(parts[1]))
			}
			for _, c := range parts[1] {
				if !strings.Contains(letterBytes, string(c)) {
					t.Errorf("GetNodeValue(%d) contains invalid character %c", tt.nodeID, c)
				}
			}
		})
	}
}

func TestCommitmentGraphConsistencyWithShuffle(t *testing.T) {
	// Create a graph with multiple nodes
	cg := coloringgraph.NewColoringGraph()
	cg.AddNode("red")
	cg.AddNode("blue")
	cg.AddNode("green")
	cg.AddEdge(0, 1)
	cg.AddEdge(1, 2)
	cg.AddEdge(2, 0)

	// Create commitment graph
	commitmentGraph := NewCommitmentGraph(cg)

	// Extract the values from the commitment graph and build a new coloring graph
	newCg := coloringgraph.NewColoringGraph()
	for i, value := range commitmentGraph.nodesValues {
		parts := strings.Split(value, "|")
		if len(parts) != 2 {
			t.Errorf("GetNodeValue(%d) returned invalid format: %s", i, value)
		}
		newCg.AddNode(coloringgraph.ColorNodeValue(parts[0]))
		newCg.Colors[parts[0]] = struct{}{}
	}
	for _, edge := range commitmentGraph.GetEdges() {
		newCg.AddEdge(edge.From, edge.To)
	}

	if !newCg.IsGraphColoringValid() {
		t.Error("newCg is not valid")
	}
}
