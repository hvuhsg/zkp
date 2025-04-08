package zkp

import (
	"crypto/sha1"
	"math/big"

	coloringgraph "github.com/hvuhsg/zkp/coloring_graph"
	commitmentgraph "github.com/hvuhsg/zkp/commitment_graph"
)

type Proofer struct {
	coloredGraph *coloringgraph.ColoringGraph
}

type CommitementGraphPayload []byte

func (cgp CommitementGraphPayload) Hash() [20]byte {
	return sha1.Sum(cgp)
}

type Proof struct {
	commitementGraphs []CommitementGraphPayload
	edgeIds           []uint64
	edgeValues        [][2]string
}

func NewProofer(coloredGraph *coloringgraph.ColoringGraph) *Proofer {
	return &Proofer{
		coloredGraph: coloredGraph,
	}
}

func (p *Proofer) CreateProof(length int) *Proof {
	commitementGraphsPayloads := make([]CommitementGraphPayload, length)
	commitementGraphs := make([]*commitmentgraph.CommitmentGraph, length)
	edgeValues := make([][2]string, length)
	edgeIds := make([]uint64, length)

	for i := range commitementGraphsPayloads {
		cg := commitmentgraph.NewCommitmentGraph(p.coloredGraph)
		commitementGraphsPayloads[i] = cg.Serialize()
		commitementGraphs[i] = cg
	}

	edgeNonce := hashModMaxUint64(commitementGraphsPayloads[len(commitementGraphsPayloads)-1].Hash())
	for i, cg := range commitementGraphs {
		edges := cg.GetEdges()

		edgeId := edgeNonce % uint64(len(edges))
		edge := edges[edgeId]
		nodev1 := cg.GetNodeValue(edge.From)
		nodev2 := cg.GetNodeValue(edge.To)

		edgeValues[i] = [2]string{nodev1, nodev2}
		edgeIds[i] = edgeId
		edgeNonce = hashModMaxUint64(commitementGraphsPayloads[i].Hash()) + edgeNonce
	}

	return &Proof{
		commitementGraphs: commitementGraphsPayloads,
		edgeValues:        edgeValues,
		edgeIds:           edgeIds,
	}
}

// HashModMaxUint64 takes a 20-byte array and returns the value modulo max uint64
func hashModMaxUint64(hash [20]byte) uint64 {
	// Convert the full 20-byte hash to a big integer
	hashInt := new(big.Int).SetBytes(hash[:])

	// Convert result to uint64 and return
	return hashInt.Uint64()
}
