package bot

import (
	"math/bits"
	"testing"
)

//func BenchmarkState_GeneratePossibleMoves(b *testing.B) {
//	state := &State{
//		Day:       8,
//		Nutrients: 20,
//		MySun:     13,
//		Trees:     [37]*Tree{},
//	}
//	// enemy trees
//	mine := false
//	state.AddTree(&Tree{CellIndex: 25, Size: SizeLarge, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 3, Size: SizeLarge, IsMine: mine, IsDormant: true})
//	state.AddTree(&Tree{CellIndex: 22, Size: SizeMedium, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 10, Size: SizeMedium, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 2, Size: SizeMedium, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 0, Size: SizeSmall, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 1, Size: SizeSeed, IsMine: mine, IsDormant: false})
//	// my trees
//	mine = true
//	state.AddTree(&Tree{CellIndex: 34, Size: SizeLarge, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 15, Size: SizeLarge, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 31, Size: SizeMedium, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 14, Size: SizeMedium, IsMine: mine, IsDormant: false})
//	state.AddTree(&Tree{CellIndex: 16, Size: SizeMedium, IsMine: mine, IsDormant: false})
//
//	for n := 0; n < b.N; n++ {
//		state.GeneratePossibleMoves()
//	}
//}

func TestState_CalculateShadows(t *testing.T) {
	state := &State{
		Trees: []Tree{36: {}},
		Num:   []int8{SizeLarge: 0},
	}
	state.AddTree(Tree{
		CellIndex: 0,
		Size:      SizeSmall,
		IsMine:    true,
		IsDormant: false,
	})

	shadows := state.Sun.CalculateShadows(state.Trees)

	// cell 1 should be shaded
	if !IsSet(shadows[0][SizeSmall], 1) {
		t.Error(IsSet(shadows[0][SizeSmall], 1))
	}
}

//
//func TestState_CalculateShadows_Complex(t *testing.T) {
//	// Day 16
//
//	state := &State{
//		Trees: []Tree{
//			1:  {CellIndex: 1, Size: SizeSeed, IsMine: true, IsDormant: true, Exists: true},
//			4:  {CellIndex: 4, Size: SizeMedium, IsMine: true, IsDormant: true, Exists: true},
//			8:  {CellIndex: 8, Size: SizeMedium, IsMine: true, IsDormant: false, Exists: true},
//			17: {CellIndex: 17, Size: SizeLarge, IsMine: true, IsDormant: false, Exists: true},
//			21: {CellIndex: 21, Size: SizeSmall, IsMine: true, IsDormant: false, Exists: true},
//			22: {CellIndex: 22, Size: SizeLarge, IsMine: true, IsDormant: false, Exists: true},
//			24: {CellIndex: 24, Size: SizeLarge, IsMine: true, IsDormant: false, Exists: true},
//			36: {CellIndex: 36, Size: SizeLarge, IsMine: true, IsDormant: false, Exists: true},
//
//			9:  {CellIndex: 9, Size: SizeLarge, IsMine: false, IsDormant: false, Exists: true},
//			15: {CellIndex: 15, Size: SizeSmall, IsMine: false, IsDormant: false, Exists: true},
//			25: {CellIndex: 25, Size: SizeLarge, IsMine: false, IsDormant: false, Exists: true},
//			27: {CellIndex: 27, Size: SizeLarge, IsMine: false, IsDormant: false, Exists: true},
//			29: {CellIndex: 29, Size: SizeLarge, IsMine: false, IsDormant: false, Exists: true},
//			31: {CellIndex: 31, Size: SizeLarge, IsMine: false, IsDormant: false, Exists: true},
//			33: {CellIndex: 31, Size: SizeSeed, IsMine: false, IsDormant: false, Exists: true},
//
//			//36: {},
//		},
//	}
//
//	shadows := Sun{Orientation: 16 % 6}.CalculateShadows(state.Trees)
//
//	expected := []int{
//		0:  3,
//		1:  2,
//		2:  3,
//		5:  3,
//		6:  2,
//		8:  1,
//		9:  3,
//		10: 0,
//		11: 3, 12: 3, 13: 3,
//		14: 2,
//		26: 3, 27: 3, 28: 3,
//		30: 2,
//		31: 1,
//		33: 3,
//		34: 3, 35: 3,
//		36: 0,
//	}
//	for i, e := range expected {
//		if int(shadows[i]) != e {
//			t.Errorf("At index %v: expected %v got %v", i, e, shadows[i])
//		}
//	}
//}

func TestState_AddTree(t *testing.T) {
	state := State{
		Trees: make([]Tree, 37),
		Num:   make([]int8, SizeLarge+1),
	}
	state.AddTree(Tree{
		CellIndex: 0,
		Size:      SizeLarge,
		IsMine:    true,
		IsDormant: false,
		Exists:    true,
	})

	if bits.OnesCount64(state.Larges) != 1 {
		t.Error(bits.OnesCount64(state.Larges))
	}
	if bits.OnesCount64(state.Smalls) != 0 {
		t.Error(bits.OnesCount64(state.Smalls))
	}

	state.AddTree(Tree{
		CellIndex: 17,
		Size:      SizeLarge,
		IsMine:    true,
		IsDormant: false,
		Exists:    true,
	})

	if bits.OnesCount64(state.Larges) != 2 {
		t.Error(bits.OnesCount64(state.Larges))
	}
	if bits.OnesCount64(state.Smalls) != 0 {
		t.Error(bits.OnesCount64(state.Smalls))
	}
}

func TestState_GatherSun(t *testing.T) {
	state := State{
		Trees: make([]Tree, 37),
		Num:   make([]int8, SizeLarge+1),
	}
	state.AddTree(Tree{
		CellIndex: 0,
		Size:      SizeLarge,
		IsMine:    true,
		IsDormant: true,
		Exists:    true,
	})

	newState := state.GatherSun()

	if newState.MySun != 3 {
		t.Error(state.MySun)
	}
}

func TestState_GeneratePossibleMoves(t *testing.T) {
	state := unmarshal("tmp.json")

	moves := state.GeneratePossibleMoves()

	t.Log(moves)
}

func TestState_DoGrow(t *testing.T) {
	state := State{
		Trees: make([]Tree, 37),
		Num:   make([]int8, SizeLarge+1),
	}
	state.AddTree(Tree{
		CellIndex: 0,
		Size:      SizeSmall,
		IsMine:    true,
		IsDormant: false,
		Exists:    true,
	})
	state.AddTree(Tree{
		CellIndex: 1,
		Size:      SizeSmall,
		IsMine:    true,
		IsDormant: false,
		Exists:    true,
	})

	newState := state.DoGrow(Action{
		Type:          Grow,
		TargetCellIdx: 0,
	})

	if bits.OnesCount64(newState.MyTrees&newState.Smalls & ^newState.DormantTrees) != 1 {
		t.Fatal(bits.OnesCount64(newState.MyTrees & newState.Smalls & newState.DormantTrees))
	}
	if bits.OnesCount64(newState.MyTrees&newState.Mediums&newState.DormantTrees) != 1 {
		t.Fatal(bits.OnesCount64(newState.MyTrees & newState.Mediums & ^newState.DormantTrees))
	}
	if !newState.Trees[0].IsDormant {
		t.Fatal("0 is not dormant")
	}
	if newState.Trees[1].IsDormant {
		t.Fatal("1 is dormant")
	}
}
