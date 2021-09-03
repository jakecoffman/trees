package game

import "sort"

type State struct {
	Board Board

	Sun       Sun
	Day       int
	Nutrients int

	// Trees: key is cell index
	Trees map[int]*Tree

	Energy  []int
	Score   []int
	Shadows [6][4]uint64
}

func New() *State {
	s := &State{
		Board:     NewBoard(),
		Sun:       Sun{},
		Day:       0,
		Nutrients: 20,
		Trees:     map[int]*Tree{},
		Energy:    []int{0, 0},
		Score:     []int{0, 0},
	}

	// this stuff can probably go in NewBoard
	for _, cell := range s.Board.Map {
		s.Board.Cells = append(s.Board.Cells, cell)
	}
	sort.Slice(s.Board.Cells, func(i, j int) bool {
		return s.Board.Cells[i].Index < s.Board.Cells[j].Index
	})
	for coord, cell := range s.Board.Map {
		cell.Neighbors = s.Board.GetNeighborIds(coord)
		s.Board.Cells[cell.Index] = cell
		s.Board.Map[coord] = cell
	}

	s.RandomizeBoard()
	s.InitStartingTrees()
	s.GatherSun()

	return s
}
