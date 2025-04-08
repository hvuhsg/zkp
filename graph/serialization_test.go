package graph

import (
	"testing"
)

func TestNodeSerialization(t *testing.T) {
	tests := []struct {
		name     string
		node     Node[IntNodeValue]
		expected []byte
	}{
		{
			name: "simple node",
			node: Node[IntNodeValue]{
				Id:    1,
				Value: 42,
			},
			expected: []byte{
				0, 2, // id size (2)
				0, 1, // id (1)
				0, 2, // value size (2)
				0, 42, // value (42)
			},
		},
		{
			name: "max values",
			node: Node[IntNodeValue]{
				Id:    0xFFFF,
				Value: 0xFFFF,
			},
			expected: []byte{
				0, 2, // id size (2)
				255, 255, // id (65535)
				0, 2, // value size (2)
				255, 255, // value (65535)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := tt.node.Serialize()
			if len(serialized) != len(tt.expected) {
				t.Errorf("serialized length = %v, want %v", len(serialized), len(tt.expected))
			}
			for i := range serialized {
				if serialized[i] != tt.expected[i] {
					t.Errorf("serialized[%d] = %v, want %v", i, serialized[i], tt.expected[i])
				}
			}

			// Test deserialization
			deserialized, _, err := DeserializeNode(serialized, DeserializeIntNodeValue)
			if err != nil {
				t.Errorf("deserialization failed: %v", err)
			}
			if deserialized.Id != tt.node.Id {
				t.Errorf("deserialized id = %v, want %v", deserialized.Id, tt.node.Id)
			}
			if deserialized.Value != tt.node.Value {
				t.Errorf("deserialized value = %v, want %v", deserialized.Value, tt.node.Value)
			}
		})
	}
}

func TestEdgeSerialization(t *testing.T) {
	tests := []struct {
		name     string
		edge     Edge
		expected []byte
	}{
		{
			name: "simple edge",
			edge: Edge{
				From: 1,
				To:   2,
			},
			expected: []byte{
				0, 2, // from size (2)
				0, 1, // from (1)
				0, 2, // to size (2)
				0, 2, // to (2)
			},
		},
		{
			name: "max values",
			edge: Edge{
				From: 0xFFFF,
				To:   0xFFFF,
			},
			expected: []byte{
				0, 2, // from size (2)
				255, 255, // from (65535)
				0, 2, // to size (2)
				255, 255, // to (65535)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := tt.edge.Serialize()
			if len(serialized) != len(tt.expected) {
				t.Errorf("serialized length = %v, want %v", len(serialized), len(tt.expected))
			}
			for i := range serialized {
				if serialized[i] != tt.expected[i] {
					t.Errorf("serialized[%d] = %v, want %v", i, serialized[i], tt.expected[i])
				}
			}

			// Test deserialization
			deserialized, err := DeserializeEdge(serialized)
			if err != nil {
				t.Errorf("deserialization failed: %v", err)
			}
			if deserialized.From != tt.edge.From {
				t.Errorf("deserialized from = %v, want %v", deserialized.From, tt.edge.From)
			}
			if deserialized.To != tt.edge.To {
				t.Errorf("deserialized to = %v, want %v", deserialized.To, tt.edge.To)
			}
		})
	}
}

func TestGraphSerialization(t *testing.T) {
	tests := []struct {
		name     string
		graph    *Graph[IntNodeValue]
		expected []byte
	}{
		{
			name:  "empty graph",
			graph: NewGraph[IntNodeValue](),
			expected: []byte{
				version,    // version
				0, 0, 0, 0, // nodes size (0)
				0, 0, 0, 0, // edges size (0)
			},
		},
		{
			name: "graph with one node",
			graph: func() *Graph[IntNodeValue] {
				g := NewGraph[IntNodeValue]()
				g.AddNode(42)
				return g
			}(),
			expected: []byte{
				version,    // version
				0, 0, 0, 8, // nodes size (8)
				0, 2, 0, 0, 0, 2, 0, 42, // node data
				0, 0, 0, 0, // edges size (0)
			},
		},
		{
			name: "graph with one edge",
			graph: func() *Graph[IntNodeValue] {
				g := NewGraph[IntNodeValue]()
				g.AddNode(42)
				g.AddNode(43)
				g.AddEdge(0, 1)
				return g
			}(),
			expected: []byte{
				version,     // version
				0, 0, 0, 16, // nodes size (16)
				0, 2, 0, 0, 0, 2, 0, 42, // node 1
				0, 2, 0, 1, 0, 2, 0, 43, // node 2
				0, 0, 0, 8, // edges size (8)
				0, 2, 0, 0, 0, 2, 0, 1, // edge data
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			serialized := tt.graph.Serialize()
			if len(serialized) != len(tt.expected) {
				t.Errorf("serialized length = %v, want %v", len(serialized), len(tt.expected))
			}
			for i := range serialized {
				if serialized[i] != tt.expected[i] {
					t.Errorf("serialized[%d] = %v, want %v", i, serialized[i], tt.expected[i])
				}
			}

			// Test deserialization
			deserialized, err := DeserializeGraph(serialized, DeserializeIntNodeValue)
			if err != nil {
				t.Errorf("deserialization failed: %v", err)
			}
			if len(deserialized.nodes) != len(tt.graph.nodes) {
				t.Errorf("deserialized nodes length = %v, want %v", len(deserialized.nodes), len(tt.graph.nodes))
			}
			if len(deserialized.edges) != len(tt.graph.edges) {
				t.Errorf("deserialized edges length = %v, want %v", len(deserialized.edges), len(tt.graph.edges))
			}
		})
	}
}
