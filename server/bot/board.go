package bot

import (
	"sort"
)

type Board struct {
	Map           map[Coord]*Cell
	Cells         []*Cell
	GoodSeedSpots [][]*Cell
	// Coordinates are ordered by Cell Index
	Coords []Coord
	Index  int8
}

const MapRingCount = 3

func NewBoard() *Board {
	b := Board{
		Map: map[Coord]*Cell{},
	}

	center := Coord{}
	b.AddCell(center, RichnessHigh)
	coord := center.Neighbor(0, 1)

	for distance := 1; distance <= MapRingCount; distance++ {
		for orientation := 0; orientation < 6; orientation++ {
			for count := 0; count < distance; count++ {
				if distance == MapRingCount {
					b.AddCell(coord, RichnessLow)
				} else if distance == MapRingCount-1 {
					b.AddCell(coord, RichnessMedium)
				} else {
					b.AddCell(coord, RichnessHigh)
				}
				coord = coord.Neighbor(int8((orientation+2)%6), 1)
			}
		}
		coord = coord.Neighbor(0, 1)
	}

	type Pair struct {
		coord Coord
		cell  *Cell
	}
	var pairs []Pair
	for coord, cell := range b.Map {
		pairs = append(pairs, Pair{coord: coord, cell: cell})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].cell.Index < pairs[j].cell.Index
	})
	for _, pair := range pairs {
		b.Coords = append(b.Coords, pair.coord)
	}

	for _, cell := range b.Map {
		b.Cells = append(b.Cells, cell)
	}
	sort.Slice(b.Cells, func(i, j int) bool {
		return b.Cells[i].Index < b.Cells[j].Index
	})
	for coord, cell := range b.Map {
		cell.Neighbors = b.GetNeighborIds(coord)
		b.Cells[cell.Index] = cell
		b.Map[coord] = cell
	}

	CalculateGoodSeedSpots(&b)

	return &b
}

func (b *Board) GetNeighborIds(coord Coord) []int8 {
	var orderedNeighborIds []int8
	for i := int8(0); i < int8(len(directions)); i++ {
		cell, ok := b.Map[coord.Neighbor(i, 1)]
		if ok {
			orderedNeighborIds = append(orderedNeighborIds, cell.Index)
		} else {
			orderedNeighborIds = append(orderedNeighborIds, -1)
		}
	}
	return orderedNeighborIds
}

func (b *Board) AddCell(coord Coord, richness int8) {
	b.Map[coord] = &Cell{
		Index:    b.Index,
		Richness: richness,
	}
	b.Index++
}

func (b *Board) Edges() []Coord {
	center := Coord{}
	var edges []Coord
	for _, coord := range b.Coords {
		if coord.DistanceTo(center) == 3 {
			edges = append(edges, coord)
		}
	}
	return edges
}

var inRangeArr = make([]Coord, 40)

func InRange(center Coord, n int8) []Coord {
	results := inRangeArr
	var cur int
	for x := -n; x <= n; x++ {
		for y := max(-n, -x-n); y <= min(+n, -x+n); y++ {
			z := -x - y
			results[cur] = center.Add(Coord{x, y, z})
			cur++
		}
	}
	return results[:cur]
}

var coordToIndex [7][7][7]int

func CoordToIndex(coord Coord) int {
	return coordToIndex[coord.x+3][coord.y+3][coord.z+3]
}

func IsCoordValid(c Coord) bool {
	return !(c.x < -3 || c.y < -3 || c.z < -3 || c.x > 3 || c.y > 3 || c.z > 3)
}

func init() {
	for i, coord := range IndexToCoord {
		coordToIndex[coord.x+3][coord.y+3][coord.z+3] = i
	}
}

var CornerIndices = []int8{19, 22, 25, 28, 31, 34}

func IsCornerIndex(i int8) bool {
	return i == 19 || i == 22 || i == 25 || i == 28 || i == 31 || i == 34
}

var IndexToCoord = []Coord{
	0:  {0, 0, 0},
	1:  {1, -1, 0},
	2:  {1, 0, -1},
	3:  {0, 1, -1},
	4:  {-1, 1, 0},
	5:  {-1, 0, 1},
	6:  {0, -1, 1},
	7:  {2, -2, 0},
	8:  {2, -1, -1},
	9:  {2, 0, -2},
	10: {1, 1, -2},
	11: {0, 2, -2},
	12: {-1, 2, -1},
	13: {-2, 2, 0},
	14: {-2, 1, 1},
	15: {-2, 0, 2},
	16: {-1, -1, 2},
	17: {0, -2, 2},
	18: {1, -2, 1},
	19: {3, -3, 0},
	20: {3, -2, -1},
	21: {3, -1, -2},
	22: {3, 0, -3},
	23: {2, 1, -3},
	24: {1, 2, -3},
	25: {0, 3, -3},
	26: {-1, 3, -2},
	27: {-2, 3, -1},
	28: {-3, 3, 0},
	29: {-3, 2, 1},
	30: {-3, 1, 2},
	31: {-3, 0, 3},
	32: {-2, -1, 3},
	33: {-1, -2, 3},
	34: {0, -3, 3},
	35: {1, -3, 2},
	36: {2, -3, 1},
}
