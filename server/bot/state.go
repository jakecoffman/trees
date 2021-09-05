package bot

import (
	"fmt"
	"log"
	"math"
	"math/bits"
	"strings"
)

type Sun struct {
	Orientation int8
}

func (s Sun) Move() Sun {
	s.Orientation = (s.Orientation + 1) % 6
	return s
}

type Trees []Tree

func (t Trees) String() string {
	var str strings.Builder
	str.WriteString("[")
	for i := 0; i < len(t); i++ {
		if t[i].Exists {
			str.WriteString(t[i].String())
			str.WriteString(" ")
		}
	}
	str.WriteString("]")
	return str.String()
}

type State struct {
	Board *Board

	Sun       Sun `json:"-"`
	Day       int8
	Nutrients int8

	// Cell index to Tree
	Trees Trees

	Seeds        uint64
	Smalls       uint64
	Mediums      uint64
	Larges       uint64
	MyTrees      uint64
	EnemyTrees   uint64
	DormantTrees uint64

	MySun, OpponentSun     int8
	MyScore, OpponentScore int
	OpponentIsWaiting      bool

	CameFrom          *State       `json:"-"`
	Shadows           [6][4]uint64 `json:"-"`
	Num               []int8       `json:"-"`
	GeneratedByAction Action       `json:"-"`
}

func (g *State) String() string {
	return fmt.Sprintf("Sun: %v Day: %v Nutrients: %v MySun: %v OpSun: %v MyScore: %v OpScore: %v Num: %v",
		g.Sun.Orientation, g.Day, g.Nutrients, g.MySun, g.OpponentSun, g.MyScore, g.OpponentScore, g.Num)
}

func (g *State) AddTree(tree Tree) {
	tree.Exists = true
	g.Trees[tree.CellIndex] = tree
	if tree.IsMine {
		g.Num[tree.Size]++
	}
	switch tree.Size {
	case SizeSeed:
		g.Seeds |= 1 << tree.CellIndex
	case SizeSmall:
		g.Smalls |= 1 << tree.CellIndex
	case SizeMedium:
		g.Mediums |= 1 << tree.CellIndex
	case SizeLarge:
		g.Larges |= 1 << tree.CellIndex
	}
	if tree.IsMine {
		g.MyTrees |= 1 << tree.CellIndex
	} else {
		g.EnemyTrees |= 1 << tree.CellIndex
	}
	if tree.IsDormant {
		g.DormantTrees |= 1 << tree.CellIndex
	} else {
		g.DormantTrees &= ^(1 << tree.CellIndex)
	}
}

type Possibilities struct {
	cursor int
	arr    []Action
}

var possibilities = Possibilities{arr: make([]Action, 100_000)}

func (p *Possibilities) Reset() {
	p.cursor = 0
}

func (p *Possibilities) Push(action Action) {
	p.arr[p.cursor] = action
	p.cursor++
}

func (p *Possibilities) Result() []Action {
	return p.arr[:p.cursor]
}

func (g *State) GeneratePossibleMoves() []Action {
	if g.Day > 23 {
		return nil
	}

	possibilities.Reset()
	possibilities.Push(Action{Type: Wait})

	//var possibleSeeds []Action
	//var possibleGrows []Action
	//var possibleCompletes []Action

	// for each tree, where they can seed
	// for each tree, if they can grow
	seedCost := g.Cost(SizeSeed)

	for i := 0; i < len(g.Trees); i++ {
		tree := g.Trees[i]
		if !tree.Exists {
			continue
		}
		if !tree.IsMine {
			continue
		}
		coord := g.Board.Coords[i]

		if g.CanSeedFrom(tree, seedCost) {
			for _, targetCoord := range InRange(coord, tree.Size) {
				if !IsCoordValid(targetCoord) {
					continue
				}
				index := CoordToIndex(targetCoord)
				targetCell := g.Board.Cells[index]
				if g.CanSeedTo(targetCell) {
					possibilities.Push(Action{
						Type:          Seed,
						SourceCellIdx: int8(i),
						TargetCellIdx: targetCell.Index,
					})
				}
			}
		}

		growCost := g.GrowthCost(tree)
		if growCost <= g.MySun && !tree.IsDormant {
			if tree.Size == SizeLarge {
				possibilities.Push(Action{Type: Complete, TargetCellIdx: int8(i)})
			} else {
				possibilities.Push(Action{Type: Grow, TargetCellIdx: int8(i)})
			}
		}
	}

	//if debug {
	//	log.Println("COMPLETE", len(possibleCompletes))
	//	log.Println("SEEDS", len(possibleSeeds))
	//	log.Println("GROWS", len(possibleGrows))
	//}

	//lines = append(lines, possibleCompletes...)
	//lines = append(lines, possibleGrows...)
	//lines = append(lines, possibleSeeds...)

	return possibilities.Result()
}

func (g *State) StartTurn() {
	g.Trees = make([]Tree, 37)
	g.CameFrom = nil
	g.Num = make([]int8, SizeLarge+1)

	g.Seeds = 0
	g.Smalls = 0
	g.Mediums = 0
	g.Larges = 0
	g.MyTrees = 0
	g.EnemyTrees = 0
	g.DormantTrees = 0
}

func (g *State) NextAction() Action {
	log.Println("DAY", g.Day, "N", g.Nutrients)

	action, _, _ := Chokudai(g, &settings)
	//action, _, _ := Jake(g, &settings)
	return action
}

var treeBaseCost = []int8{0, 1, 3, 7}

const (
	CostLifecycleEnd    = 4
	RichnessBonusMedium = 2
	RichnessBonusHigh   = 4
)

// Cost returns the cost to build a tree of size treeSize
func (g *State) Cost(treeSize int8) int8 {
	return treeBaseCost[treeSize] + g.Num[treeSize]
	//for i := 0; i < len(g.Trees); i++ {
	//	tree := g.Trees[i]
	//	if !tree.Exists {
	//		continue
	//	}
	//	if tree.Size == treeSize && tree.IsMine {
	//		cost++
	//	}
	//}
	//return cost
}

func (g *State) CanSeedFrom(tree Tree, seedCost int8) bool {
	return seedCost <= g.MySun && tree.Size > SizeSeed && !tree.IsDormant
}

func (g *State) CanSeedTo(target Cell) bool {
	if target.Richness == RichnessUnusable {
		return false
	}
	tree := g.Trees[target.Index]
	return !tree.Exists
}

// GrowthCost returns the cost to grow tree by 1
func (g *State) GrowthCost(tree Tree) int8 {
	targetSize := tree.Size + 1
	if targetSize > SizeLarge {
		return CostLifecycleEnd
	}
	return g.Cost(targetSize)
}

// Clone clones the game state
func (g *State) Clone() *State {
	//start := time.Now()

	state := &State{
		Board: g.Board, // board is immutable
		Trees: make([]Tree, 37),
		Num:   make([]int8, SizeLarge+1),
	}

	//if time.Now().Sub(start) > 6 * time.Millisecond {
	//	log.Println("TOOK TOO LONG:", time.Now().Sub(start))
	//}

	copy(state.Trees, g.Trees)

	//if time.Now().Sub(start) > 6 * time.Millisecond {
	//	log.Println("TOOK TOO LONG:", time.Now().Sub(start))
	//}

	state.Sun.Orientation = g.Sun.Orientation
	state.Day = g.Day
	state.Nutrients = g.Nutrients
	state.MySun = g.MySun
	state.MyScore = g.MyScore
	//state.Trees = newTrees
	//state.Shadows = make([]int8, 37)
	//copy(state.Shadows, g.Shadows)
	state.Shadows = g.Shadows // shadows only change when WAIT
	state.CameFrom = nil

	state.MyTrees = g.MyTrees
	state.EnemyTrees = g.EnemyTrees
	state.Seeds = g.Seeds
	state.Smalls = g.Smalls
	state.Mediums = g.Mediums
	state.Larges = g.Larges
	state.DormantTrees = g.DormantTrees

	copy(state.Num, g.Num)
	return state
}

func (g *State) DoGrow(action Action) *State {
	newState := g.Clone()

	//coord := board.Coords[action.TargetCellIdx]
	cell := g.Board.Cells[action.TargetCellIdx]
	tree := newState.Trees[cell.Index]
	//if tree == nil {
	//	panic("Tree not found")
	//}
	//if !tree.IsMine {
	//	panic("Tree not mine")
	//}
	//if tree.IsDormant {
	//	panic("Tree is dormant")
	//}
	//if tree.Size >= SizeLarge {
	//	panic("Tree too large")
	//}
	cost := newState.GrowthCost(tree)
	sun := newState.MySun
	//if sun < cost {
	//	panic("Can't afford")
	//}
	newState.Num[tree.Size]--
	newState.Num[tree.Size+1]++

	switch tree.Size {
	case SizeSeed:
		newState.Seeds &= ^(1 << cell.Index)
		newState.Smalls |= 1 << cell.Index
	case SizeSmall:
		newState.Smalls &= ^(1 << cell.Index)
		newState.Mediums |= 1 << cell.Index
	case SizeMedium:
		newState.Mediums &= ^(1 << cell.Index)
		newState.Larges |= 1 << cell.Index
	}

	newState.MySun = sun - cost
	newState.Trees[cell.Index].Size++
	newState.Trees[cell.Index].IsDormant = true
	newState.DormantTrees |= 1 << cell.Index

	return newState
}

func (g *State) DoComplete(action Action) *State {
	newState := g.Clone()

	//coord := board.Coords[action.TargetCellIdx]
	cell := g.Board.Cells[action.TargetCellIdx]
	tree := newState.Trees[cell.Index]
	//if tree == nil {
	//	panic("Tree not found")
	//}
	//if !tree.IsMine {
	//	panic("Tree not mine")
	//}
	//if tree.IsDormant {
	//	panic("Tree is dormant")
	//}
	//if tree.Size < SizeLarge {
	//	panic("Tree too small")
	//}
	cost := newState.GrowthCost(tree)
	sun := newState.MySun
	//if sun < cost {
	//	panic("Can't afford")
	//}

	newState.MySun = sun - cost

	points := int(newState.Nutrients)
	if cell.Richness == RichnessMedium {
		points += RichnessBonusMedium
	} else if cell.Richness == RichnessHigh {
		points += RichnessBonusHigh
	}
	newState.MyScore += points
	newState.Trees[cell.Index].Exists = false
	newState.Nutrients = max(0, newState.Nutrients-1)
	//newState.Shadows = newState.Sun.CalculateShadows(newState.Trees)
	newState.Num[SizeLarge]--
	newState.Larges &= ^(1 << cell.Index)
	newState.MyTrees &= ^(1 << cell.Index)

	return newState
}

func (g *State) DoSeed(action Action) *State {
	newState := g.Clone()

	//targetCoord := board.Coords[action.TargetCellIdx]
	//sourceCoord := board.Coords[action.SourceCellIdx]

	targetCell := g.Board.Cells[action.TargetCellIdx]
	sourceCell := g.Board.Cells[action.SourceCellIdx]

	//if _, ok := g.Trees[targetCell.Index]; ok {
	//	panic("Cell is not empty")
	//}
	sourceTree := &newState.Trees[sourceCell.Index]
	//if sourceTree == nil {
	//	panic("Tree not found")
	//}
	//if sourceTree.Size == SizeSeed {
	//	panic("Seeds can't seed")
	//}
	//if !sourceTree.IsMine {
	//	panic("Tree is not mine")
	//}
	//if sourceTree.IsDormant {
	//	panic("Tree is dormant")
	//}

	//distance := sourceCoord.DistanceTo(targetCoord)
	//if distance > sourceTree.Size {
	//	panic("Tree can't seed that far")
	//}
	//if targetCell.Richness == RichnessUnusable {
	//	panic("Cell is unusable")
	//}

	costOfSeed := newState.Cost(SizeSeed)
	//sun := newState.MySun
	//if sun < costOfSeed {
	//	panic("Not enough sun to seed")
	//}
	sourceTree.IsDormant = true
	newState.DormantTrees |= 1 << sourceTree.CellIndex
	newState.MySun -= costOfSeed
	newState.AddTree(Tree{
		CellIndex: targetCell.Index,
		Size:      SizeSeed,
		IsMine:    true,
		IsDormant: true,
	})

	return newState
}

func (g *State) GatherSun() *State {
	if g.Day == 24 {
		log.Panic("TIME WARP ", g.Day)
	}
	newState := g.Clone()
	newState.Sun = newState.Sun.Move()
	newState.Shadows = newState.Sun.CalculateShadows(g.Board, newState.Trees)

	smallSun := bits.OnesCount64(g.MyTrees & g.Smalls & ^(g.Shadows[0][SizeSmall]))
	mediumSun := bits.OnesCount64(g.MyTrees & g.Mediums & ^(g.Shadows[0][SizeMedium]))
	largeSun := bits.OnesCount64(g.MyTrees & g.Larges & ^(g.Shadows[0][SizeLarge]))
	newState.MySun += int8(smallSun + mediumSun*2 + largeSun*3)

	for i := 0; i < len(g.Trees); i++ {
		newState.Trees[i].IsDormant = false
	}
	newState.DormantTrees = 0
	newState.Day++

	return newState
}

func (g *State) EndScore() int {
	return g.MyScore + int(math.Floor(float64(g.MySun)/3))
}
