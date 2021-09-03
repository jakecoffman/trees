package game

type Sun struct {
	Orientation int
}

func (s Sun) Move() Sun {
	s.Orientation = (s.Orientation + 1) % 6
	return s
}

// CalculateShadows returns shadows for the current sun.
func (s Sun) CalculateShadows(state *State) [SizeLarge + 1][]int {
	shadows := [4][]int{
		// doing initialization to make frontend code simpler
		0: {},
		1: {},
		2: {},
		3: {},
	}
	for j := range state.Trees {
		tree := state.Trees[j]
		coord := state.Board.Coords[tree.CellIndex]
		for i := 1; i <= tree.Size; i++ { // huge bug here <=
			tempCoord := coord.Neighbor(s.Orientation, i)
			if IsCoordValid(tempCoord) {
				shadows[tree.Size] = append(shadows[tree.Size], CoordToIndex(tempCoord))
			}
		}
	}

	return shadows
}
