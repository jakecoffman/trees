package bot

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/bits"
	"sort"
	"testing"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func unmarshal(fileName string) (state *State) {
	d, err := ioutil.ReadFile("scenarios/" + fileName)
	check(err)
	err = json.Unmarshal(d, &state)
	check(err)
	state.Board = NewBoard()
	state.Sun = Sun{Orientation: state.Day % 6}
	state.Shadows = state.Sun.CalculateShadows(state.Board, state.Trees)
	state.CameFrom = nil
	state.GeneratedByAction = Action{}
	state.Num = make([]int8, SizeLarge+1)
	for i := 0; i < 37; i++ {
		t := state.Trees[i]
		if t.Exists && t.IsMine {
			state.Num[t.Size]++
		}
	}
	return
}

func TestChokudai_Turn1(t *testing.T) {
	state := unmarshal("turn1.json")

	start := time.Now()
	action, _, _ := Chokudai(state, &settings)
	t.Log("Took", time.Now().Sub(start))

	if action.Type != Wait {
		t.Error(action)
	}
}

func TestChokudai_Turn2(t *testing.T) {
	state := unmarshal("turn2.json")
	log.Println("Starting state is", state)

	start := time.Now()
	_, _, path := Chokudai(state, &settings)
	t.Log("Took", time.Now().Sub(start))

	expected := []int8{Grow, Wait, Grow, Seed, Wait}
	log.Println("Expected final score", path[len(path)-1].Priority(&settings))
	for i := range expected {
		if path[i].GeneratedByAction.Type != expected[i] {
			t.Error(i, path[i].GeneratedByAction, "!=", Action{Type: expected[i]})
		}
	}
}

//func TestChokudai_Turn14(t *testing.T) {
//	state := unmarshal("turn14.json")
//	cp := make([]Cell, len(Cells))
//	copy(cp, Cells)
//	defer copy(Cells, cp)
//	for _, i := range []int{1, 3, 4, 6, 9, 15, 22, 31} {
//		Cells[i].Richness = RichnessUnusable
//	}
//	log.Println("Starting state is", state)
//
//	start := time.Now()
//	Chokudai(state, &settings)
//	t.Log("Took", time.Now().Sub(start))
//}

func TestStuff(t *testing.T) {
	// https://www.codingame.com/replay/558505027
	state := unmarshal("turn2.json")

	//1
	{
		var possibleStates []*State
		moves := state.GeneratePossibleMoves()
		for _, m := range moves {
			if m.Type == Seed && state.Trees[m.SourceCellIdx].Size == SizeSmall {
				continue
			}
			possibleStates = append(possibleStates, applyMove(m, state))
		}
		sort.Slice(possibleStates, func(i, j int) bool {
			return possibleStates[i].Priority(&settings) > possibleStates[j].Priority(&settings)
		})
		//for _, s := range possibleStates {
		//	t.Log(s.Priority(&settings), s.GeneratedByAction)
		//}
		state = possibleStates[0]
	}
	if state.GeneratedByAction.Type != Grow && state.Priority(&settings) != 5 {
		t.Fatal(state.GeneratedByAction)
	}
	if state.MySun != 1 {
		t.Fatal(state.MySun)
	}

	//2
	{
		var possibleStates []*State
		moves := state.GeneratePossibleMoves()
		for _, m := range moves {
			if m.Type == Seed && state.Trees[m.SourceCellIdx].Size == SizeSmall {
				continue
			}
			possibleStates = append(possibleStates, applyMove(m, state))
		}
		sort.Slice(possibleStates, func(i, j int) bool {
			return possibleStates[i].Priority(&settings) > possibleStates[j].Priority(&settings)
		})
		//for _, s := range possibleStates {
		//	t.Log(s.Priority(&settings), s.GeneratedByAction)
		//}
		state = possibleStates[0]
	}
	if state.GeneratedByAction.Type != Wait {
		t.Fatal(state.GeneratedByAction)
	}
	if state.MySun != 4 {
		t.Fatal(state.MySun, bits.OnesCount64(state.Mediums), "M", bits.OnesCount64(state.Smalls), "S")
	}

	//3
	{
		var possibleStates []*State
		moves := state.GeneratePossibleMoves()
		for _, m := range moves {
			if m.Type == Seed && state.Trees[m.SourceCellIdx].Size == SizeSmall {
				continue
			}
			possibleStates = append(possibleStates, applyMove(m, state))
		}
		//t.Log(possibleStates)
		sort.Slice(possibleStates, func(i, j int) bool {
			return possibleStates[i].Priority(&settings) > possibleStates[j].Priority(&settings)
		})
		//for _, s := range possibleStates {
		//	t.Log(s.Priority(&settings), s.GeneratedByAction)
		//}
		state = possibleStates[0]
	}
	if state.GeneratedByAction.Type != Grow {
		t.Fatal(state.GeneratedByAction)
	}

	//4
	{
		var possibleStates []*State
		moves := state.GeneratePossibleMoves()
		for _, m := range moves {
			if m.Type == Seed && state.Trees[m.SourceCellIdx].Size == SizeSmall {
				continue
			}
			possibleStates = append(possibleStates, applyMove(m, state))
		}
		//t.Log(possibleStates)
		sort.Slice(possibleStates, func(i, j int) bool {
			return possibleStates[i].Priority(&settings) > possibleStates[j].Priority(&settings)
		})
		//for _, s := range possibleStates {
		//	t.Log(s.Priority(&settings), s.GeneratedByAction)
		//}
		state = possibleStates[0]
	}
	if state.GeneratedByAction.Type != Seed {
		t.Fatal(state.GeneratedByAction)
	}

	//5
	{
		var possibleStates []*State
		moves := state.GeneratePossibleMoves()
		for _, m := range moves {
			if m.Type == Seed && state.Trees[m.SourceCellIdx].Size == SizeSmall {
				continue
			}
			possibleStates = append(possibleStates, applyMove(m, state))
		}
		//t.Log(possibleStates)
		sort.Slice(possibleStates, func(i, j int) bool {
			return possibleStates[i].Priority(&settings) > possibleStates[j].Priority(&settings)
		})
		//for _, s := range possibleStates {
		//	t.Log(s.Priority(&settings), s.GeneratedByAction)
		//}
		state = possibleStates[0]
	}
	if state.GeneratedByAction.Type != Wait {
		t.Fatal(state.GeneratedByAction)
	}

	//6
	{
		var possibleStates []*State
		moves := state.GeneratePossibleMoves()
		for _, m := range moves {
			if m.Type == Seed && state.Trees[m.SourceCellIdx].Size == SizeSmall {
				continue
			}
			possibleStates = append(possibleStates, applyMove(m, state))
		}
		//t.Log(possibleStates)
		sort.Slice(possibleStates, func(i, j int) bool {
			return possibleStates[i].Priority(&settings) > possibleStates[j].Priority(&settings)
		})
		//for _, s := range possibleStates {
		//	t.Log(s.Priority(&settings), s.GeneratedByAction)
		//}
		state = possibleStates[0]
	}
	log.Println("FINAL SCORE", state.Priority(&settings))
	if state.GeneratedByAction.Type != Grow {
		t.Fatal(state.GeneratedByAction)
	}
}

func TestChokudai_No_Move_Detected(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	state := unmarshal("tmp.json")

	start := time.Now()
	action, _, path := Chokudai(state, &settings)
	t.Log("Took", time.Now().Sub(start))

	if action.Type != Complete {
		t.Error(action)
	}
	for _, p := range path {
		t.Log(p.Priority(&settings), "SUN:", p.MySun, p.GeneratedByAction)
	}
}

func TestChokudai_Stupid_Seed(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	state := unmarshal("dumb_seed.json")

	action, _, _ := Chokudai(state, &settings)

	if action.Type == Seed {
		t.Error(action)
	}
}

func TestChokudai_FinalRound_Complete(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	state := &State{
		Sun:       Sun{},
		Day:       23,
		Nutrients: 0,
		Trees: []Tree{
			0: {
				CellIndex: 0,
				Size:      SizeLarge,
				IsMine:    true,
				IsDormant: false,
				Exists:    true,
			},
			36: {},
		},
		MyScore: 0,
		Num:     []int8{SizeLarge: 0},
	}
	state.MySun = state.Cost(3)
	//log.Println(len(state.GeneratePossibleMoves()))

	action, _, _ := Chokudai(state, &settings)

	if action.Type != Complete {
		t.Error(action.Type)
	}
}

//func TestChokudai_Seed_Maybe(t *testing.T) {
//	log.SetFlags(log.Lshortfile)
//
//	state := &State{
//		Sun:       Sun{},
//		Day:       0,
//		Nutrients: 0,
//		Trees: []Tree{
//			0: {
//				CellIndex: 0,
//				Size:      SizeLarge,
//				IsMine:    true,
//				IsDormant: false,
//				Exists:    true,
//			},
//			36: {},
//		},
//		MyScore: 0,
//		Num:     []int8{SizeLarge: 0},
//	}
//	state.MySun = state.Cost(0)
//	//log.Println(state.GeneratePossibleMoves())
//
//	action, _, path := Chokudai(state, &settings)
//
//	//log.Println(len(path))
//	//for _, s := range path {
//	//	log.Println(s.GeneratedByAction)
//	//}
//
//	if action.Type != Seed {
//		t.Error(action.Type, path)
//	}
//}

func TestChokudai_Completey_Boyz(t *testing.T) {
	log.SetFlags(log.Lshortfile)

	state := &State{
		Sun:       Sun{},
		Day:       23,
		Nutrients: 0,
		MyScore:   0,
		Trees:     []Tree{36: {}},
		Num:       []int8{SizeLarge: 0},
	}
	state.AddTree(Tree{
		CellIndex: 0,
		Size:      SizeLarge,
		IsMine:    true,
		IsDormant: false,
	})
	state.AddTree(Tree{
		CellIndex: 1,
		Size:      SizeLarge,
		IsMine:    true,
		IsDormant: false,
	})
	state.MySun = state.Cost(3) * 2
	//state.GeneratePossibleMoves(true)

	action, _, _ := Chokudai(state, &settings)

	if action.Type != Complete {
		t.Error(action.Type)
	}
}

func TestChokudai_Start(t *testing.T) {
	state := &State{
		Sun:       Sun{},
		Day:       0,
		Nutrients: 20,
		Trees:     []Tree{36: {}},
		MySun:     0,
		MyScore:   0,
		Num:       []int8{SizeLarge: 0},
	}
	state.AddTree(Tree{
		CellIndex: 27,
		Size:      SizeSmall,
		IsMine:    true,
		IsDormant: false,
	})
	state.AddTree(Tree{
		CellIndex: 21,
		Size:      SizeSmall,
		IsMine:    true,
		IsDormant: false,
	})

	action, _, _ := Chokudai(state, &settings)

	if action.Type != Wait {
		t.Error(action.Type)
	}
}

//func TestChokudai_Best(t *testing.T) {
//	if testing.Short() {
//		t.Skip()
//	}
//
//	var bestScore int
//	var bestSetting Settings
//	rand.Seed(time.Now().Unix())
//
//	debug.SetGCPercent(-1)
//
//	// simulate games
//	for i := 0; i < 100; i++ {
//		runtime.GC()
//		// randomize settings
//		newSettings := settings.Tweak()
//
//		// starting state always the same
//		state := &State{
//			Sun:       Sun{},
//			Day:       0,
//			Nutrients: 20,
//			Trees:     []Tree{36: {}},
//			MySun:     0,
//			MyScore:   0,
//			Num:       []int8{SizeLarge: 0},
//		}
//		state.AddTree(Tree{
//			CellIndex: 27,
//			Size:      SizeSmall,
//			IsMine:    true,
//			IsDormant: false,
//		})
//		state.AddTree(Tree{
//			CellIndex: 21,
//			Size:      SizeSmall,
//			IsMine:    true,
//			IsDormant: false,
//		})
//		state.Shadows = state.Sun.Move().CalculateShadows(state.Trees)
//
//		var turns int
//		for state.Day < 23 {
//			start := time.Now()
//
//			action, _, _ := Chokudai(state, &newSettings)
//			turns++
//
//			diff := time.Now().Sub(start)
//			if diff > 91*time.Millisecond {
//				log.Println(time.Now().Sub(start), state.Day)
//			}
//			if diff >= 100*time.Millisecond {
//				log.Println("TIMED OUT:", time.Now().Sub(start))
//				break
//			}
//
//			//if best == nil {
//			//log.Println("NO RESULTS:", settings)
//			//continue
//			//}
//			var newState *State
//
//			switch action.Type {
//			case Wait:
//				newState = state.GatherSun()
//				//log.Println("WAIT", t, current.State.MyScore, "->", nextState.MyScore)
//			case Seed:
//				newState = state.DoSeed(action)
//				//log.Println("SEED", t, current.State.MyScore, "->", nextState.MyScore)
//			case Complete:
//				newState = state.DoComplete(action)
//				//log.Println("COMP", t, current.State.MyScore, "->", nextState.MyScore)
//			case Grow:
//				newState = state.DoGrow(action)
//				//log.Println("SEED", t, current.State.MyScore, "->", nextState.MyScore)
//			default:
//				log.Panic("Invalid state: ", action.Type)
//			}
//
//			state.Day = newState.Day
//			state.Sun = newState.Sun
//			state.Nutrients = newState.Nutrients
//			state.MySun = newState.MySun
//			state.MyScore = newState.MyScore
//			copy(state.Trees, newState.Trees)
//			state.Shadows = newState.Shadows
//		}
//		numTrees := 0
//		for _, t := range state.Trees {
//			if t.Exists && t.IsMine {
//				numTrees++
//			}
//		}
//		log.Println("GAME", i, "SCORE", state.MyScore, "TURNS", turns, "LEFT SUN", state.MySun, "TREES", numTrees)
//		if state.MyScore > bestScore {
//			bestScore = state.MyScore
//			bestSetting = newSettings
//			log.Println("NEW HIGH SCORE!", newSettings)
//		}
//	}
//	t.Log("Best:", bestScore, bestSetting)
//}
//
//// go test -cpuprofile cpu.prof -memprofile mem.prof -bench .
//func BenchmarkChokudai(b *testing.B) {
//	{
//		f, err := os.Create("cpu.prof")
//		if err != nil {
//			log.Fatal("could not create CPU profile: ", err)
//		}
//		defer f.Close() // error handling omitted for example
//		if err := pprof.StartCPUProfile(f); err != nil {
//			log.Fatal("could not start CPU profile: ", err)
//		}
//		defer pprof.StopCPUProfile()
//	}
//	//{
//	//	f, err := os.Create("mem.prof")
//	//	if err != nil {
//	//		log.Fatal("could not create memory profile: ", err)
//	//	}
//	//	defer f.Close() // error handling omitted for example
//	//	runtime.GC()    // get up-to-date statistics
//	//	if err := pprof.WriteHeapProfile(f); err != nil {
//	//		log.Fatal("could not write memory profile: ", err)
//	//	}
//	//}
//
//	state := &State{
//		Sun:               Sun{Orientation: 3},
//		Day:               9,
//		Nutrients:         19,
//		MySun:             3,
//		OpponentSun:       6,
//		MyScore:           0,
//		OpponentScore:     20,
//		OpponentIsWaiting: true,
//		Trees:             []Tree{36: {}},
//		Num:               []int8{SizeLarge: 0},
//	}
//	state.AddTree(Tree{CellIndex: 0, Size: 2, IsMine: true, IsDormant: true})
//	state.AddTree(Tree{CellIndex: 1, Size: 1, IsMine: true, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 2, Size: 1, IsMine: true, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 3, Size: 2, IsMine: false, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 4, Size: 0, IsMine: true, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 5, Size: 3, IsMine: true, IsDormant: true})
//	state.AddTree(Tree{CellIndex: 6, Size: 2, IsMine: true, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 15, Size: 2, IsMine: false, IsDormant: true})
//	state.AddTree(Tree{CellIndex: 22, Size: 2, IsMine: true, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 31, Size: 2, IsMine: false, IsDormant: false})
//	state.AddTree(Tree{CellIndex: 35, Size: 2, IsMine: true, IsDormant: false})
//
//	for i := 0; i < b.N; i++ {
//		Chokudai(state, &settings)
//	}
//}
//
//func BenchmarkState_Clone(b *testing.B) {
//	state := &State{
//		Sun:       Sun{},
//		Day:       0,
//		Nutrients: 0,
//		Trees: []Tree{
//			27: {
//				CellIndex: 27,
//				Size:      SizeSmall,
//				IsMine:    true,
//				IsDormant: false,
//			},
//			21: {
//				CellIndex: 21,
//				Size:      SizeSmall,
//				IsMine:    true,
//				IsDormant: false,
//			},
//			36: {},
//		},
//		MySun:   0,
//		MyScore: 0,
//	}
//
//	for n := 0; n < b.N; n++ {
//		state.Clone()
//		poolCursor = 0
//		//treePoolCursor = 0
//	}
//}
