package commitmentgraph

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"

	coloringgraph "github.com/hvuhsg/zkp/coloring_graph"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateRandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

type CommitmentGraph struct {
	*coloringgraph.ColoringGraph
	nodesValues []string
}

func NewCommitmentGraph(cg *coloringgraph.ColoringGraph) *CommitmentGraph {
	cg.ShuffleColors()

	nodesValues := make([]string, len(cg.GetNodes()))
	for i, node := range cg.GetNodes() {
		nodeStringValue := string(node.Value)
		randomString := generateRandomString(10)
		nodeValue := nodeStringValue + "|" + randomString
		hash := sha1.Sum([]byte(nodeValue))
		node.Value = coloringgraph.ColorNodeValue(hex.EncodeToString(hash[:]))
		nodesValues[i] = nodeValue
	}

	return &CommitmentGraph{
		ColoringGraph: cg,
		nodesValues:   nodesValues,
	}
}

func (cg *CommitmentGraph) GetNodeValue(id int) string {
	return cg.nodesValues[id]
}
