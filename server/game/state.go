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
	s = s.GatherSun()

	return s
}

func (s *State) Clone() *State {
	newState := &State{
		Trees: map[int]*Tree{},
	}

	for k, v := range s.Trees {
		tree := *v
		newState.Trees[k] = &tree
		if newState.Trees[k] == v {
			panic("PROGRAMMER WRONG")
		}
	}

	// Board never changes so we don't copy it
	newState.Board = s.Board
	newState.Sun = s.Sun
	newState.Day = s.Day
	newState.Nutrients = s.Nutrients
	newState.Energy = []int{s.Energy[0], s.Energy[1]}
	newState.Score = []int{s.Score[0], s.Score[1]}

	return newState
}
