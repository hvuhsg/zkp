# Zero-Knowledge Proof Implementation

This project implements a zero-knowledge proof system for graph coloring problems using Go. It provides a way to prove knowledge of a valid graph coloring without revealing the actual coloring solution.

## Overview

The implementation consists of three main components:
1. **Proofer**: Creates zero-knowledge proofs for graph coloring solutions
2. **Verifier**: Validates the proofs without revealing the actual coloring
3. **Graph Structures**: Implements various graph types and coloring mechanisms

## Features

- Zero-knowledge proof generation for graph coloring
- Cryptographic commitment scheme using SHA-1
- Graph serialization and deserialization
- Comprehensive test suite
- Benchmarking capabilities

## Project Structure

```
.
├── coloring_graph/     # Graph coloring implementation
├── commitment_graph/   # Commitment scheme for proofs
├── graph/             # Base graph data structures
├── proofer.go         # Proof generation
├── verifier.go        # Proof verification
├── randomizer.go      # Random number generation for proofs
└── *_test.go          # Test files
```

## Requirements

- Go 1.24.1 or higher
- Dependencies:
  - github.com/stretchr/testify v1.10.0

## Installation

1. Clone the repository:
```bash
git clone https://github.com/hvuhsg/zkp.git
cd zkp
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### Creating a Proof

```go
// Create a colored graph
coloredGraph := coloringgraph.NewColoringGraph(...)

// Initialize the proofer
proofer := zkp.NewProofer(coloredGraph)

// Generate a proof
proof := proofer.CreateProof(length)
```

### Verifying a Proof

```go
// Verify a proof
isValid := proof.Verify()
```

## Testing

Run the test suite:
```bash
go test ./...
```

Run benchmarks:
```bash
go test -bench=. ./...
```

## Security

This implementation uses SHA-1 for cryptographic commitments. While this is sufficient for demonstration purposes, for production use, consider using a more modern cryptographic hash function.
