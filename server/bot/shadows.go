package bot

func SetBit(v uint64, index int8) uint64 {
	v |= 1 << index
	return v
}

func ClearBit(v uint64, index int8) uint64 {
	v &= ^(1 << index)
	return v
}

func IsSet(v uint64, index int8) bool {
	return (v & 1 << index) > 0
}

// CalculateShadows -- 0 indicates this turn's shadows
func (s Sun) CalculateShadows(b *Board, trees []Tree) [6][SizeLarge + 1]uint64 {
	var shadows [6][4]uint64
	for o := 0; o < 6; o++ {
		for j := 0; j < len(trees); j++ {
			tree := trees[j]
			if !tree.Exists {
				continue
			}
			coord := b.Coords[tree.CellIndex]
			for i := int8(1); i <= tree.Size; i++ { // huge bug here <=
				tempCoord := coord.Neighbor(s.Orientation, i)
				if IsCoordValid(tempCoord) {
					shadows[o][tree.Size] |= 1 << CoordToIndex(tempCoord)
				}
			}
		}
		s = s.Move()
	}

	return shadows
}
