package graph

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

// Serialize the node into a byte array
// the format is as follows:
// [id_size][id][value_size][value]
func (n *Node[T]) Serialize() []byte {
	var buf bytes.Buffer

	// Write ID size (2 bytes for uint16)
	idSize := make([]byte, 2)
	binary.BigEndian.PutUint16(idSize, 2)
	buf.Write(idSize)

	// Write ID value
	idBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(idBytes, n.Id)
	buf.Write(idBytes)

	// Write value size (2 bytes for uint16)
	valueBytes := n.Value.Serialize()
	valueSize := make([]byte, 2)

	if len(valueBytes) > math.MaxUint16 {
		panic("value size too large")
	}

	binary.BigEndian.PutUint16(valueSize, uint16(len(valueBytes)))
	buf.Write(valueSize)

	// Write value
	buf.Write(valueBytes)

	return buf.Bytes()
}

// Serialize the edge into a byte array
// the format is as follows:
// [from_size][from_value][to_size][to_value]
func (e *Edge) Serialize() []byte {
	var buf bytes.Buffer

	// Write from size (2 bytes for uint16)
	fromSize := make([]byte, 2)
	binary.BigEndian.PutUint16(fromSize, 2)
	buf.Write(fromSize)

	// Write from value
	fromBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(fromBytes, uint16(e.From))
	buf.Write(fromBytes)

	// Write to size (2 bytes for uint16)
	toSize := make([]byte, 2)
	binary.BigEndian.PutUint16(toSize, 2)
	buf.Write(toSize)

	// Write to value
	toBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(toBytes, uint16(e.To))
	buf.Write(toBytes)

	return buf.Bytes()
}

// Serialize the graph into a byte array
// the format is as follows:
// [version][nodes_size][node1_size][node1_value][node2_size][node2_value]...[edges_size][edge1_from_size][edge1_from_value][edge1_to_size][edge1_to_value]...
func (g *Graph[T]) Serialize() []byte {
	var buf bytes.Buffer

	// Write version (1 byte)
	buf.Write([]byte{version})

	// Serialize nodes
	var nodesBuffer bytes.Buffer
	for _, node := range g.nodes {
		nodesBuffer.Write(node.Serialize())
	}

	// Write nodes size (4 bytes for uint32 since it could be large)
	nodesSize := make([]byte, 4)
	binary.BigEndian.PutUint32(nodesSize, uint32(nodesBuffer.Len()))
	buf.Write(nodesSize)
	buf.Write(nodesBuffer.Bytes())

	// Serialize edges
	var edgesBuffer bytes.Buffer
	for _, edge := range g.edges {
		edgesBuffer.Write(edge.Serialize())
	}

	// Write edges size (4 bytes for uint32 since it could be large)
	edgesSize := make([]byte, 4)
	binary.BigEndian.PutUint32(edgesSize, uint32(edgesBuffer.Len()))
	buf.Write(edgesSize)
	buf.Write(edgesBuffer.Bytes())

	return buf.Bytes()
}

// DeserializeNode creates a Node from a byte array
// the format is as follows:
// [id_size][id][value_size][value]
func DeserializeNode[T NodeValue](data []byte, valueDeserializer func([]byte) (T, error)) (*Node[T], uint, error) {
	totalSize := uint(0)

	if len(data) < 8 { // Minimum size for a node (2+2+2+2)
		return nil, 0, fmt.Errorf("data too short for node")
	}

	// Read ID size
	idSize := binary.BigEndian.Uint16(data[0:2])
	if idSize != 2 {
		return nil, 0, fmt.Errorf("invalid id size: %d", idSize)
	}
	totalSize += 2

	// Read ID value
	id := binary.BigEndian.Uint16(data[2:4])
	totalSize += 2

	// Read value size
	valueSize := binary.BigEndian.Uint16(data[4:6])
	totalSize += 2
	// Read value
	valueBytes := data[6 : 6+valueSize]
	value, err := valueDeserializer(valueBytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to deserialize node value: %w", err)
	}
	totalSize += uint(len(valueBytes))
	return &Node[T]{
		Id:    id,
		Value: T(value),
	}, totalSize, nil
}

// DeserializeEdge creates an Edge from a byte array
// the format is as follows:
// [from_size][from_value][to_size][to_value]
func DeserializeEdge(data []byte) (*Edge, error) {
	if len(data) < 8 { // Minimum size for an edge (2+2+2+2)
		return nil, fmt.Errorf("data too short for edge")
	}

	// Read from size
	fromSize := binary.BigEndian.Uint16(data[0:2])
	if fromSize != 2 {
		return nil, fmt.Errorf("invalid from size: %d", fromSize)
	}

	// Read from value
	from := int(binary.BigEndian.Uint16(data[2:4]))

	// Read to size
	toSize := binary.BigEndian.Uint16(data[4:6])
	if toSize != 2 {
		return nil, fmt.Errorf("invalid to size: %d", toSize)
	}

	// Read to value
	to := int(binary.BigEndian.Uint16(data[6:8]))

	return &Edge{
		From: from,
		To:   to,
	}, nil
}

// DeserializeGraph creates a Graph from a byte array
// the format is as follows:
// [version][nodes_size][node1_size][node1_value][node2_size][node2_value]...[edges_size][edge1_from_size][edge1_from_value][edge1_to_size][edge1_to_value]...
func DeserializeGraph[T NodeValue](data []byte, valueDeserializer func([]byte) (T, error)) (*Graph[T], error) {
	if len(data) < 9 { // Minimum size for a graph (1+4+4)
		return nil, fmt.Errorf("data too short for graph")
	}

	// Read version
	if data[0] != version {
		return nil, fmt.Errorf("invalid version: %d", data[0])
	}

	// Read nodes size
	nodesSize := binary.BigEndian.Uint32(data[1:5])
	if len(data) < int(5+nodesSize) {
		return nil, fmt.Errorf("data too short for nodes")
	}

	// Deserialize nodes
	nodesData := data[5 : 5+nodesSize]
	nodes := make([]Node[T], 0)
	offset := 0
	for offset < len(nodesData) {
		node, totalSize, err := DeserializeNode(nodesData[offset:], valueDeserializer)
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize node: %w", err)
		}
		nodes = append(nodes, *node)
		offset += int(totalSize)
	}

	// Read edges size
	edgesOffset := 5 + nodesSize
	if len(data) < int(edgesOffset+4) {
		return nil, fmt.Errorf("data too short for edges size")
	}
	edgesSize := binary.BigEndian.Uint32(data[edgesOffset : edgesOffset+4])

	// Deserialize edges
	if len(data) < int(edgesOffset+4+edgesSize) {
		return nil, fmt.Errorf("data too short for edges")
	}
	edgesData := data[edgesOffset+4 : edgesOffset+4+edgesSize]
	edges := make([]Edge, 0)
	offset = 0
	for offset < len(edgesData) {
		edge, err := DeserializeEdge(edgesData[offset:])
		if err != nil {
			return nil, fmt.Errorf("failed to deserialize edge: %w", err)
		}
		edges = append(edges, *edge)
		offset += 8 // Each edge is 8 bytes
	}

	return &Graph[T]{
		nodes: nodes,
		edges: edges,
	}, nil
}
