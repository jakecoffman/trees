package game

type Sun struct {
	Orientation int
}

func (s Sun) Move() Sun {
	s.Orientation = (s.Orientation + 1) % 6
	return s
}

func (s Sun) CalculateShadows(state *State) [6][SizeLarge + 1]uint64 {
	var shadows [6][4]uint64
	nextSun := s.Move()
	for o := 0; o < 6; o++ {
		for j := range state.Trees {
			tree := state.Trees[j]
			coord := state.Board.Coords[tree.CellIndex]
			for i := 1; i <= tree.Size; i++ { // huge bug here <=
				tempCoord := coord.Neighbor(nextSun.Orientation, i)
				if IsCoordValid(tempCoord) {
					shadows[o][tree.Size] |= 1 << CoordToIndex(tempCoord)
				}
			}
		}
		nextSun = nextSun.Move()
	}

	return shadows
}
