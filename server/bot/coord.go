package bot

import "fmt"

var directions = [][]int8{{1, -1, 0}, {+1, 0, -1}, {0, +1, -1}, {-1, +1, 0}, {-1, 0, +1}, {0, -1, +1}}

type Coord struct {
	x, y, z int8
}

func CalculateGoodSeedSpots(b *Board) {
	// compute good seed spots
	compute := func(c Coord) []Coord {
		return []Coord{
			c.Add(Coord{+1, +1, -2}),
			c.Add(Coord{-1, +2, -1}),
			c.Add(Coord{-2, +1, +1}),
			c.Add(Coord{-1, -1, +2}),
			c.Add(Coord{+1, -2, +1}),
			c.Add(Coord{+2, -1, -1}),

			c.Add(Coord{+1, +2, -3}),
			c.Add(Coord{+2, +1, -3}),
			c.Add(Coord{-1, +3, -2}),
			c.Add(Coord{-2, +3, -1}),
			c.Add(Coord{-3, +2, +1}),
			c.Add(Coord{-3, +1, +2}),
			c.Add(Coord{-2, -1, +3}),
			c.Add(Coord{-1, -2, +3}),
			c.Add(Coord{+1, -3, +2}),
			c.Add(Coord{+2, -3, +1}),
			c.Add(Coord{+3, -2, -1}),
			c.Add(Coord{+3, -1, -2}),
		}
	}
	for _, cell := range b.Cells {
		var spots []*Cell
		coord := b.Coords[cell.Index]

		possibles := compute(coord)
		for _, c := range possibles {
			// if it's on the board
			if cell, ok := b.Map[c]; ok {
				spots = append(spots, cell)
			}
		}

		b.GoodSeedSpots = append(b.GoodSeedSpots, spots)
	}
}

func (c Coord) Add(other Coord) Coord {
	return Coord{c.x + other.x, c.y + other.y, c.z + other.z}
}

func (c Coord) Neighbor(orientation, distance int8) Coord {
	nx := c.x + directions[orientation][0]*distance
	ny := c.y + directions[orientation][1]*distance
	nz := c.z + directions[orientation][2]*distance
	return Coord{nx, ny, nz}
}

func (c Coord) DistanceTo(destination Coord) int8 {
	return (abs(c.x-destination.x) + abs(c.y-destination.y) + abs(c.z-destination.z)) / 2
}

func (c Coord) String() string {
	return fmt.Sprintf("[%d,%d,%d]", c.x, c.y, c.z)
}

func (c Coord) Opposite() Coord {
	return Coord{-c.x, -c.y, -c.z}
}

func abs(a int8) int8 {
	if a < 0 {
		return -a
	}
	return a
}
