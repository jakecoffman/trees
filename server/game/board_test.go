package game

import (
	"log"
	"testing"
)

func TestBoard_Richness(t *testing.T) {
	state := New()

	if state.Board.Map[Coord{}].Richness != RichnessHigh {
		t.Errorf("WRONG")
	}
	if state.Board.Map[IndexToCoord[25]].Richness != RichnessLow {
		t.Errorf("LOW")
	}
	if state.Board.Cells[0].Richness != RichnessHigh {
		t.Errorf("WRONG")
	}
	if state.Board.Cells[25].Richness != RichnessLow {
		t.Errorf("LOW")
	}
}

func TestBoard_Neighbors(t *testing.T) {
	lookup := map[int][]int{}
	state := New()
	for _, cell := range state.Board.Cells {
		lookup[cell.Index] = cell.Neighbors
	}
	log.Printf("%#v", lookup)
}
