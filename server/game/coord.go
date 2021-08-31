package game

import "fmt"

var directions = [][]int{{1, -1, 0}, {+1, 0, -1}, {0, +1, -1}, {-1, +1, 0}, {-1, 0, +1}, {0, -1, +1}}

type Coord struct {
	x, y, z int
}

func (c Coord) Add(other Coord) Coord {
	return Coord{c.x + other.x, c.y + other.y, c.z + other.z}
}

func (c Coord) Neighbor(orientation, distance int) Coord {
	nx := c.x + directions[orientation][0]*distance
	ny := c.y + directions[orientation][1]*distance
	nz := c.z + directions[orientation][2]*distance
	return Coord{nx, ny, nz}
}

func (c Coord) DistanceTo(destination Coord) int {
	return (abs(c.x-destination.x) + abs(c.y-destination.y) + abs(c.z-destination.z)) / 2
}

func (c Coord) String() string {
	return fmt.Sprintf("[%d,%d,%d]", c.x, c.y, c.z)
}

func (c Coord) Opposite() Coord {
	return Coord{-c.x, -c.y, -c.z}
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
