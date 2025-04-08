package coloringgraph

import "math/rand/v2"

func (cg *ColoringGraph) ShuffleColors() {
	colorsSlice := make([]string, 0, len(cg.Colors))
	for color := range cg.Colors {
		colorsSlice = append(colorsSlice, color)
	}

	// Create a new shuffled slice
	shuffledSlice := make([]string, len(colorsSlice))
	copy(shuffledSlice, colorsSlice)
	rand.Shuffle(len(shuffledSlice), func(i, j int) {
		shuffledSlice[i], shuffledSlice[j] = shuffledSlice[j], shuffledSlice[i]
	})

	shuffleMap := make(map[string]string)
	for i, color := range colorsSlice {
		shuffleMap[color] = shuffledSlice[i]
	}

	for _, node := range cg.GetNodes() {
		node.Value = ColorNodeValue(shuffleMap[string(node.Value)])
	}
}
