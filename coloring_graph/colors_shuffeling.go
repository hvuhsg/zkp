package coloringgraph

import "math/rand/v2"

func (cg *ColoringGraph) ShuffleColors() {
	// Get all unique colors from nodes
	colorsMap := make(map[string]struct{})
	for _, node := range cg.GetNodes() {
		colorsMap[string(node.Value)] = struct{}{}
	}

	// Convert colors map to slice
	colorsSlice := make([]string, 0, len(colorsMap))
	for color := range colorsMap {
		colorsSlice = append(colorsSlice, color)
	}

	// Create a new shuffled slice
	shuffledSlice := make([]string, len(colorsSlice))
	copy(shuffledSlice, colorsSlice)
	rand.Shuffle(len(shuffledSlice), func(i, j int) {
		shuffledSlice[i], shuffledSlice[j] = shuffledSlice[j], shuffledSlice[i]
	})

	// Create mapping from original colors to shuffled colors
	shuffleMap := make(map[string]string)
	for i, color := range colorsSlice {
		shuffleMap[color] = shuffledSlice[i]
	}

	// Apply the shuffle to all nodes
	for _, node := range cg.GetNodes() {
		node.Value = ColorNodeValue(shuffleMap[string(node.Value)])
	}
}
