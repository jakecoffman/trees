package main

import (
	"log"
	"reflect"
	"testing"
)

func TestClone(t *testing.T) {
	state := &State{MyScore: 1, Trees: []Tree{36: {}}}
	newState := state.Clone()

	if newState.MyScore != 1 {
		t.Error(newState.MyScore)
	}

	state = &State{
		MyScore: 2,
		Trees: []Tree{
			36: {
				CellIndex: 36,
				Size:      SizeLarge,
				IsMine:    true,
				IsDormant: true,
				Exists:    true,
			},
		},
		Num: []int8{SizeLarge: 0},
	}

	newerState := state.Clone()

	if len(newerState.Trees) != len(state.Trees) {
		t.Fatal("LEARN TO COUNT")
	}

	if !reflect.DeepEqual(newerState, state) {
		t.Error("NOT DEEP EQUAL")
	}
	for i := 0; i < 37; i++ {
		if newerState.Trees[i] != state.Trees[i] {
			t.Errorf("Tree %v not equal: %#v != %#v", i, newerState.Trees[i], state.Trees[i])
		}
	}
}

func TestBoard_Richness(t *testing.T) {
	if board.Map[Coord{}].Richness != RichnessHigh {
		t.Errorf("WRONG")
	}
	if board.Map[IndexToCoord[25]].Richness != RichnessLow {
		t.Errorf("LOW")
	}
	if Cells[0].Richness != RichnessHigh {
		t.Errorf("WRONG")
	}
	if Cells[25].Richness != RichnessLow {
		t.Errorf("LOW")
	}
}

func TestBoard_Neighbors(t *testing.T) {
	lookup := map[int8][]int8{}
	for _, cell := range Cells {
		lookup[cell.Index] = cell.Neighbors
	}
	log.Printf("%#v", lookup)
}
