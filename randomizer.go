package zkp

import (
	"math/big"
	"math/rand"
)

func commitmentsHash(commitmentsPayloads []CommitementGraphPayload) [20]byte {
	allpayloads := make([]byte, 0)
	for _, cp := range commitmentsPayloads {
		allpayloads = append(allpayloads, cp...)
	}

	return CommitementGraphPayload(allpayloads).Hash()
}

// HashModMaxUint64 takes a 20-byte array and returns the value modulo max uint64
func hashModMaxUint64(hash [20]byte) uint64 {
	// Convert the full 20-byte hash to a big integer
	hashInt := new(big.Int).SetBytes(hash[:])

	// Convert result to uint64 and return
	return hashInt.Uint64()
}

func NewRandomizerFromCommitments(commitments []CommitementGraphPayload) *rand.Rand {
	allCommitmentsHash := commitmentsHash(commitments)
	return rand.New(rand.NewSource(int64(hashModMaxUint64(allCommitmentsHash))))
}
