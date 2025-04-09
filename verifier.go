package zkp

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	coloredgraph "github.com/hvuhsg/zkp/coloring_graph"
	graph "github.com/hvuhsg/zkp/graph"
)

func (p *Proof) Verify() bool {
	// Deserialize the commitment graphs
	graphs := make([]*graph.Graph[coloredgraph.ColorNodeValue], len(p.commitementGraphs))
	for i, cg := range p.commitementGraphs {
		var err error
		graphs[i], err = graph.DeserializeGraph(cg, coloredgraph.DeserializeColorNodeValue)
		if err != nil {
			return false
		}
	}

	randomizer := NewRandomizerFromCommitments(p.commitementGraphs)

	// Verify the proof
	for i, g := range graphs {
		edgeNonce := randomizer.Uint64()

		// Verify edge values are not the same
		if !isEdgeValuesValid(p.edgeValues[i]) {
			return false
		}

		edges := g.GetEdges()

		// Verify the edge id is valid
		if !isExpectedEdgeIdValid(len(edges), edgeNonce, p.edgeIds[i]) {
			return false
		}

		nodes := g.GetNodes()
		edge := edges[p.edgeIds[i]]
		node1 := nodes[edge.From]
		node2 := nodes[edge.To]

		edgeValue1Hash := sha1.Sum([]byte(p.edgeValues[i][0]))
		edgeValue2Hash := sha1.Sum([]byte(p.edgeValues[i][1]))
		edgeValue1HashString := hex.EncodeToString(edgeValue1Hash[:])
		edgeValue2HashString := hex.EncodeToString(edgeValue2Hash[:])

		// verify hash of edge values is the same as nodes values
		if string(node1.Value) != edgeValue1HashString || string(node2.Value) != edgeValue2HashString {
			return false
		}
	}
	return true
}

func isEdgeValuesValid(edgeValues [2]string) bool {
	color1, err := getColorFromNodeValue(edgeValues[0])
	if err != nil {
		return false
	}
	color2, err := getColorFromNodeValue(edgeValues[1])
	if err != nil {
		return false
	}
	if color1 == color2 {
		return false
	}
	return true
}

func isExpectedEdgeIdValid(edgesCount int, edgeNonce uint64, edgeId uint64) bool {
	expectedEdgeId := edgeNonce % uint64(edgesCount)
	return expectedEdgeId == edgeId
}

func getColorFromNodeValue(nodeValue string) (string, error) {
	parts := strings.Split(nodeValue, "|")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid node value: %s", nodeValue)
	}
	return parts[0], nil
}
