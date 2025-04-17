package soduku

import coloringgraph "github.com/hvuhsg/zkp/coloring_graph"

type SodukuBoard struct {
	board          [][]uint
	visableIndexes [][2]uint
}

func NewSodukuBoard(initialBoard [][]uint) SodukuBoard {
	visableIndexs := make([][2]uint, 0)

	for row := range initialBoard {
		for col := range row {
			if initialBoard[row][col] != 0 {
				visableIndexs = append(visableIndexs, [2]uint{uint(row), uint(col)})
			}
		}
	}

	return SodukuBoard{
		board:          initialBoard,
		visableIndexes: visableIndexs,
	}
}

func (sb SodukuBoard) ConvertToColoringGraph() *coloringgraph.ColoringGraph {
	graph := coloringgraph.NewColoringGraph()

	// Add "colors"
	for i := range 9 {
		color := string(i + 1)
		graph.Colors[color] = struct{}{}
		graph.AddNode(coloringgraph.ColorNodeValue(color))
	}

	// Add all nodes
	for row := range sb.board {
		for col := range row {
			color := string(sb.board[row][col])
			graph.AddNode(coloringgraph.ColorNodeValue(color))
		}
	}

	// Add initial constraints
	for _, visableIndexs := range sb.visableIndexes {
		value := sb.board[visableIndexs[0]][visableIndexs[1]]
		color := string(value)

		// Add edges
		for i := range 9 {
			c := string(i + 1)
			if c == color {
				continue
			}
			visableNodeIndex := visableIndexs[0]*9 + visableIndexs[1] + 9
			graph.AddEdge(int(visableNodeIndex), i)
		}
	}

	// Add all constraints
	// TODO: add all soduku constraints

	return graph
}
