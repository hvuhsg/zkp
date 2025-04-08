package zkp

import (
	"testing"

	coloringgraph "github.com/hvuhsg/zkp/coloring_graph"
	"github.com/stretchr/testify/assert"
)

func TestVerifyValidProof(t *testing.T) {
	// Create a test graph
	graph := coloringgraph.NewColoringGraph()
	graph.AddNode(coloringgraph.ColorNodeValue("red"))
	graph.AddNode(coloringgraph.ColorNodeValue("blue"))
	graph.AddNode(coloringgraph.ColorNodeValue("green"))
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(0, 2)

	// Create a proof
	proofer := NewProofer(graph)
	proof := proofer.CreateProof(3)

	// Test that a valid proof verifies
	assert.True(t, proof.Verify(), "Valid proof should verify successfully")

	// Test that modifying edge values makes verification fail
	originalEdgeValue := proof.edgeValues[0]
	proof.edgeValues[0] = [2]string{"red|abc123", "red|def456"} // Same color
	assert.False(t, proof.Verify(), "Proof with same colors should fail verification")
	proof.edgeValues[0] = originalEdgeValue // Restore original value

	// Test that modifying edge IDs makes verification fail
	originalEdgeId := proof.edgeIds[0]
	proof.edgeIds[0] = 999 // Invalid edge ID
	assert.False(t, proof.Verify(), "Proof with invalid edge ID should fail verification")
	proof.edgeIds[0] = originalEdgeId // Restore original value

	// Test that modifying commitment graphs makes verification fail
	originalCommitmentGraph := proof.commitementGraphs[0]
	proof.commitementGraphs[0] = []byte("invalid")
	assert.False(t, proof.Verify(), "Proof with invalid commitment graph should fail verification")
	proof.commitementGraphs[0] = originalCommitmentGraph // Restore original value
}

func TestVerifyInvalidProof(t *testing.T) {
	// Create a test graph
	graph := coloringgraph.NewColoringGraph()
	graph.AddNode(coloringgraph.ColorNodeValue("red"))
	graph.AddNode(coloringgraph.ColorNodeValue("blue"))
	graph.AddNode(coloringgraph.ColorNodeValue("blue"))
	graph.AddEdge(0, 1)
	graph.AddEdge(1, 2)
	graph.AddEdge(0, 2)

	// Create a proof
	proofer := NewProofer(graph)
	proof := proofer.CreateProof(100)

	// Test that a valid proof verifies
	assert.False(t, proof.Verify(), "Invalid proof should fail verification")
}

func TestIsEdgeValuesValid(t *testing.T) {
	// Test valid edge values with different colors
	assert.True(t, isEdgeValuesValid([2]string{"red|abc123", "blue|def456"}))
	assert.True(t, isEdgeValuesValid([2]string{"blue|abc123", "green|def456"}))
	assert.True(t, isEdgeValuesValid([2]string{"green|abc123", "red|def456"}))

	// Test invalid edge values with same color
	assert.False(t, isEdgeValuesValid([2]string{"red|abc123", "red|def456"}))
	assert.False(t, isEdgeValuesValid([2]string{"blue|abc123", "blue|def456"}))
	assert.False(t, isEdgeValuesValid([2]string{"green|abc123", "green|def456"}))

	// Test invalid edge values with wrong format
	assert.False(t, isEdgeValuesValid([2]string{"red", "blue"}))
	assert.False(t, isEdgeValuesValid([2]string{"red|abc123", "blue"}))
	assert.False(t, isEdgeValuesValid([2]string{"", ""}))
}

func TestIsExpectedEdgeIdValid(t *testing.T) {
	// Test valid edge IDs
	assert.True(t, isExpectedEdgeIdValid(5, 3, 3))
	assert.True(t, isExpectedEdgeIdValid(5, 8, 3))
	assert.True(t, isExpectedEdgeIdValid(5, 13, 3))

	// Test invalid edge IDs
	assert.False(t, isExpectedEdgeIdValid(5, 3, 4))
	assert.False(t, isExpectedEdgeIdValid(5, 8, 2))
	assert.False(t, isExpectedEdgeIdValid(5, 13, 1))
}

func TestGetColorFromNodeValue(t *testing.T) {
	// Test valid node values
	color, err := getColorFromNodeValue("red|abc123")
	assert.NoError(t, err)
	assert.Equal(t, "red", color)

	color, err = getColorFromNodeValue("blue|def456")
	assert.NoError(t, err)
	assert.Equal(t, "blue", color)

	// Test invalid node values
	_, err = getColorFromNodeValue("red")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid node value")

	_, err = getColorFromNodeValue("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid node value")

	_, err = getColorFromNodeValue("red|abc123|extra")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid node value")
}
