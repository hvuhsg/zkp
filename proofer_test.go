package zkp

import (
	"strings"
	"testing"

	coloringgraph "github.com/hvuhsg/zkp/coloring_graph"
	"github.com/stretchr/testify/assert"
)

func TestNewProofer(t *testing.T) {
	// Create a test graph
	graph := coloringgraph.NewColoringGraph()
	graph.AddNode(coloringgraph.ColorNodeValue("red"))
	graph.AddNode(coloringgraph.ColorNodeValue("blue"))
	graph.AddNode(coloringgraph.ColorNodeValue("green"))
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(0, 2)

	// Test creating a new proofer
	proofer := NewProofer(graph)
	assert.NotNil(t, proofer)
	assert.Equal(t, graph, proofer.coloredGraph)
}

func TestCreateProof(t *testing.T) {
	// Create a test graph
	graph := coloringgraph.NewColoringGraph()
	graph.AddNode(coloringgraph.ColorNodeValue("red"))
	graph.AddNode(coloringgraph.ColorNodeValue("blue"))
	graph.AddNode(coloringgraph.ColorNodeValue("green"))
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(0, 2)

	proofer := NewProofer(graph)
	proof := proofer.CreateProof(3)

	// Verify proof structure
	assert.NotNil(t, proof)
	assert.Len(t, proof.commitementGraphs, 3)
	assert.Len(t, proof.edgeValues, 3)
	assert.Len(t, proof.edgeIds, 3)

	// Verify edge values are valid colors
	validColors := map[string]bool{"red": true, "blue": true, "green": true}
	for _, edgeValue := range proof.edgeValues {
		color1 := strings.Split(edgeValue[0], "|")[0]
		color2 := strings.Split(edgeValue[1], "|")[0]
		assert.True(t, validColors[color1], "Edge value should be a valid color")
		assert.True(t, validColors[color2], "Edge value should be a valid color")
	}
}

func TestHashModMaxUint64(t *testing.T) {
	// Test with a known hash value
	testHash := [20]byte{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18, 19, 20,
	}
	result := hashModMaxUint64(testHash)
	assert.NotZero(t, result, "Hash should not be zero")
}

func TestCommitmentGraphPayloadHash(t *testing.T) {
	// Test with a simple payload
	payload := CommitementGraphPayload([]byte("test"))
	hash := payload.Hash()
	assert.NotZero(t, hash, "Hash should not be zero")

	// Test that same payload produces same hash
	hash2 := payload.Hash()
	assert.Equal(t, hash, hash2, "Same payload should produce same hash")

	// Test that different payloads produce different hashes
	payload2 := CommitementGraphPayload([]byte("test2"))
	hash3 := payload2.Hash()
	assert.NotEqual(t, hash, hash3, "Different payloads should produce different hashes")
}
