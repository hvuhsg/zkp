package zkp

import (
	"testing"

	coloringgraph "github.com/hvuhsg/zkp/coloring_graph"
)

func createCircularGraph(nodesCount int) *coloringgraph.ColoringGraph {
	graph := coloringgraph.NewColoringGraph()

	colors := []string{"red", "blue", "green", "yellow", "purple", "orange", "pink", "brown", "gray", "black"}

	for i := 0; i < nodesCount; i++ {
		graph.AddNode(coloringgraph.ColorNodeValue(colors[i%len(colors)]))
	}

	for i := 0; i < nodesCount; i++ {
		graph.AddEdge(i, (i+1)%nodesCount)
	}
	return graph
}

func BenchmarkProofCreation(b *testing.B) {
	// Create a test graph with 10 nodes
	graph := createCircularGraph(10000)

	proofer := NewProofer(graph)

	b.ResetTimer()
	for b.Loop() {
		proofer.CreateProof(10)
	}
}

func BenchmarkProofVerification(b *testing.B) {
	// Create a test graph with 10 nodes
	graph := createCircularGraph(100)

	proofer := NewProofer(graph)
	proof := proofer.CreateProof(100)

	b.ResetTimer()
	for b.Loop() {
		proof.Verify()
	}
}

func BenchmarkProofCreationLargeGraph(b *testing.B) {
	// Create a test graph with 100 nodes
	graph := createCircularGraph(10000)

	proofer := NewProofer(graph)

	b.ResetTimer()
	for b.Loop() {
		proofer.CreateProof(10000)
	}
}

func BenchmarkProofVerificationLargeGraph(b *testing.B) {
	// Create a test graph with 1000 nodes
	graph := createCircularGraph(10000)

	proofer := NewProofer(graph)
	proof := proofer.CreateProof(10000)

	b.ResetTimer()
	for b.Loop() {
		proof.Verify()
	}
}
