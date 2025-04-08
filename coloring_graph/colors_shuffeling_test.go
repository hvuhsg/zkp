package coloringgraph

import (
	"testing"
)

func TestShuffleColors(t *testing.T) {
	// Create a test graph with known colors
	cg := NewColoringGraph()
	cg.AddNode(ColorNodeValue("red"))
	cg.AddNode(ColorNodeValue("blue"))
	cg.AddNode(ColorNodeValue("green"))
	cg.AddNode(ColorNodeValue("yellow"))
	cg.AddEdge(0, 1)
	cg.AddEdge(1, 2)
	cg.AddEdge(2, 3)

	// Store original colors
	originalColors := make(map[int]string)
	for i, node := range cg.GetNodes() {
		originalColors[i] = string(node.Value)
	}

	// Shuffle colors
	cg.ShuffleColors()

	// Verify that all nodes still have colors
	for i, node := range cg.GetNodes() {
		if node.Value == "" {
			t.Errorf("Node %d has no color after shuffling", i)
		}
	}

	// Verify that the graph structure remains intact
	if len(cg.GetNodes()) != 4 {
		t.Errorf("Expected 4 nodes, got %d", len(cg.GetNodes()))
	}

	// Verify that all edges are preserved
	edges := cg.GetEdges()
	if len(edges) != 3 {
		t.Errorf("Expected 3 edges, got %d", len(edges))
	}

	// Verify that no two adjacent nodes have the same color
	for _, edge := range edges {
		sourceNode := cg.GetNodes()[edge.From]
		targetNode := cg.GetNodes()[edge.To]
		if sourceNode.Value == targetNode.Value {
			t.Errorf("Adjacent nodes %d and %d have the same color: %s",
				edge.From, edge.To, sourceNode.Value)
		}
	}

	// Verify that the color mapping is bijective (one-to-one)
	colorCount := make(map[string]int)
	for _, node := range cg.GetNodes() {
		colorCount[string(node.Value)]++
	}
	for color, count := range colorCount {
		if count != 1 {
			t.Errorf("Color %s appears %d times, expected 1", color, count)
		}
	}

	// Verify that all original colors are still present
	shuffledColors := make(map[string]bool)
	for _, node := range cg.GetNodes() {
		shuffledColors[string(node.Value)] = true
	}
	for _, originalColor := range originalColors {
		if !shuffledColors[originalColor] {
			t.Errorf("Original color %s is missing after shuffling", originalColor)
		}
	}
}
