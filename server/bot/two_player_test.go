package bot

import (
	"fmt"
	"log"
	"math/rand"
	"runtime/debug"
	"testing"
	"time"
)

func TestTwoPlayers(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	var bestScore int
	//var bestSetting Settings
	rand.Seed(time.Now().Unix())

	debug.SetGCPercent(-1)

	//limit = 10 * time.Millisecond
	settings := Settings{
		MinimalEconomy: 16,
		MaximumEconomy: 25,
	}
	var mySettings = settings
	var oppSettings = settings

	// There are 2 ways to use this:
	// 1. Set this to false and manually give opponent different settings.
	// 2. Set this to true and each game opponent will change settings until it wins.
	// Once it wins it will set mySettings to the winning value. That way it finds
	// values that are better (in theory).
	const shouldTweakSettings = false
	oppSettings.InDev = true

	var wins [2]int

	// simulate games
	for i := 0; i < 100; i++ {
		state := &State{
			Board:     NewBoard(),
			Nutrients: 20,
			Trees:     make([]Tree, 37),
			Num:       make([]int8, SizeLarge+1),
		}
		RandomizeBoard(state)
		InitStartingTrees(state)
		state.Shadows = state.Sun.CalculateShadows(state.Board, state.Trees)
		state = GatherSun(state)
		var err error

		var myLastAction, oppLastAction Action
		for state.Day < 23 {

			var myAction Action
			if myLastAction.Type == Wait && oppLastAction.Type != Wait {
				myAction = myLastAction
			} else {
				myAction = nextAction(state, mySettings)
			}

			var oppAction Action
			if oppLastAction.Type == Wait && myLastAction.Type != Wait {
				oppAction = oppLastAction
			} else {
				swapSides(state)
				oppAction = nextAction(state, oppSettings)
				swapSides(state)
			}

			log.Println(myAction, oppAction)

			myLastAction, oppLastAction = myAction, oppAction

			if myAction.Type == Seed && oppAction.Type == Seed {
				if myAction.TargetCellIdx == oppAction.TargetCellIdx {
					// seed is not planted, refund sun, make source trees dormant
					state.Trees[myAction.SourceCellIdx].IsDormant = true
					state.Trees[oppAction.SourceCellIdx].IsDormant = true
					continue
				}
			}

			if myAction.Type == Wait && oppAction.Type == Wait {
				state.Day++
				state.Sun = state.Sun.Move()
				state = GatherSun(state)
				continue
			}
			if myAction.Type == Seed {
				state, err = ApplySeed(state, &myAction)
				if err != nil {
					panic(err)
				}
			}
			if oppAction.Type == Seed {
				swapSides(state)
				state, err = ApplySeed(state, &oppAction)
				if err != nil {
					panic(err)
				}
				swapSides(state)
			}
			var nutrientsLost int8
			if myAction.Type == Complete {
				nutrientsLost++
				state, err = ApplyComplete(state, &myAction)
				if err != nil {
					panic(err)
				}
			}
			if oppAction.Type == Complete {
				nutrientsLost++
				swapSides(state)
				state, err = ApplyComplete(state, &oppAction)
				if err != nil {
					panic(err)
				}
				swapSides(state)
			}
			state.Nutrients = max(0, state.Nutrients-nutrientsLost)
			if myAction.Type == Grow {
				state, err = ApplyGrow(state, &myAction)
				if err != nil {
					panic(err)
				}
			}
			if oppAction.Type == Grow {
				swapSides(state)
				state, err = ApplyGrow(state, &oppAction)
				if err != nil {
					panic(err)
				}
				swapSides(state)
			}
		}

		if state.MyScore > bestScore {
			bestScore = state.MyScore
			//log.Println("HIGH SCORE", bestScore, mySettings)
		}
		if state.OpponentScore > bestScore {
			bestScore = state.OpponentScore
			log.Println("OPP HIGH SCORE", bestScore, oppSettings)
		}

		if state.MyScore == state.OpponentScore {
			log.Println(i, "TIE", "SCORE", state.MyScore, "-", state.OpponentScore)
		} else if state.MyScore > state.OpponentScore {
			log.Println(i, "MY WIN", "SCORE", state.MyScore, "-", state.OpponentScore)
			if shouldTweakSettings {
				oppSettings = settings.Tweak()
			}
			wins[0]++
		} else {
			log.Println(i, "OP WIN", "SCORE", state.MyScore, "-", state.OpponentScore)
			if shouldTweakSettings {
				settings = oppSettings
				mySettings = settings
				oppSettings = settings.Tweak()
				// mySettings become best settings, then tweak opponent from there
			}
			wins[1]++
		}

		t.Logf("Wins: %v - %v", wins[0], wins[1])
	}
	t.Log("Best:", bestScore, settings)
}

func InitStartingTrees(s *State) {
	startingCoords := s.Board.Edges()

	// remove if richness is null
	{
		var dedupeStartingCoords []Coord
		for _, c := range startingCoords {
			if s.Board.Map[c].Richness != RichnessUnusable {
				dedupeStartingCoords = append(dedupeStartingCoords, c)
			}
		}
		startingCoords = dedupeStartingCoords
	}

	// while size < 2, try init starting trees
	var validCoords []Coord
	const numPlayers = 2
	for len(validCoords) < startingTreeCount*numPlayers {
		validCoords = append(validCoords, tryInitStartingTrees(startingCoords)...)
	}

	for i := 0; i < startingTreeCount; i++ {
		{
			cell := s.Board.Map[validCoords[2*i]]
			s.AddTree(Tree{
				CellIndex: cell.Index,
				Size:      SizeSmall,
				IsMine:    true,
				IsDormant: false,
			})
		}
		{
			cell := s.Board.Map[validCoords[2*i+1]]
			s.AddTree(Tree{
				CellIndex: cell.Index,
				Size:      SizeSmall,
				IsMine:    false,
				IsDormant: false,
			})
		}
	}
}

const (
	startingTreeCount    = 2
	startingTreeDistance = 2
)

func tryInitStartingTrees(startingCoords []Coord) []Coord {
	var coords []Coord

	var availableCoords []Coord
	for _, coord := range startingCoords {
		availableCoords = append(availableCoords, coord)
	}

	for i := 0; i < startingTreeCount; i++ {
		if len(availableCoords) == 0 {
			return coords
		}
		r := rand.Intn(len(availableCoords))
		normalCoord := availableCoords[r]
		oppositeCoord := normalCoord.Opposite()
		for j := 0; j < len(availableCoords); {
			coord := availableCoords[j]
			if coord.DistanceTo(normalCoord) <= startingTreeDistance || coord.DistanceTo(oppositeCoord) <= startingTreeDistance {
				availableCoords = append(availableCoords[:j], availableCoords[j+1:]...)
			} else {
				j++
			}
		}
		coords = append(coords, normalCoord)
		coords = append(coords, oppositeCoord)
	}

	return coords
}

const maxEmptyCells = 10

// RandomizeBoard adds unused cells
func RandomizeBoard(s *State) {
	// reset all cells to the proper richness
	for coord, cell := range s.Board.Map {
		cell.Neighbors = s.Board.GetNeighborIds(coord)
		s.Board.Cells[cell.Index] = cell
		s.Board.Map[coord] = cell
	}

	wantedEmptyCells := rand.Intn(maxEmptyCells + 1)
	actuallyEmptyCells := 0
	for actuallyEmptyCells < wantedEmptyCells-1 {
		index := rand.Intn(len(s.Board.Map))
		randCoord := s.Board.Coords[index]
		if s.Board.Map[randCoord].Richness != RichnessUnusable {
			s.Board.Map[randCoord].Richness = RichnessUnusable
			actuallyEmptyCells++
			if randCoord != randCoord.Opposite() {
				s.Board.Map[randCoord.Opposite()].Richness = RichnessUnusable
				actuallyEmptyCells++
			}
		}
	}
}

// GatherSun calculates shadows and adds energy to players
func GatherSun(s *State) *State {
	state := s.Clone()
	state.Shadows = state.Sun.CalculateShadows(s.Board, state.Trees)
	for i := range state.Trees {
		tree := state.Trees[i]
		state.Trees[i].IsDormant = false
		if tree.IsMine {
			if !IsSet(state.Shadows[0][tree.Size], tree.CellIndex) {
				state.MySun += tree.Size
			}
		} else {
			if !IsSet(state.Shadows[0][tree.Size], tree.CellIndex) {
				state.OpponentSun += tree.Size
			}
		}
	}
	return state
}

func nextAction(state *State, s Settings) Action {
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

	//log.Println(state)
	//log.Println(state.Trees)

	start := time.Now()

	myAction, _, _ := Chokudai(state, &s)

	diff := time.Now().Sub(start)
	if diff > 97*time.Millisecond {
		log.Println(time.Now().Sub(start), state.Day)
	}
	if diff >= 100*time.Millisecond {
		log.Println("TIMED OUT:", time.Now().Sub(start))
		myAction = Action{Type: Wait}
	}
	return myAction
}

func swapSides(state *State) {
	state.MySun, state.OpponentSun = state.OpponentSun, state.MySun
	state.MyScore, state.OpponentScore = state.OpponentScore, state.MyScore
	state.MyTrees, state.EnemyTrees = state.EnemyTrees, state.MyTrees
	for i := 0; i < 37; i++ {
		t := &state.Trees[i]
		t.IsMine = !t.IsMine
	}
}

// ApplyGrow makes player grow a tree
func ApplyGrow(s *State, action *Action) (*State, error) {
	state := s.Clone()

	cell := state.Board.Cells[action.TargetCellIdx]
	tree := state.Trees[cell.Index]
	if !tree.Exists {
		return nil, fmt.Errorf("tree not found at %v", cell.Index)
	}
	if !tree.IsMine {
		return nil, fmt.Errorf("tree at %v not owned by player", tree.CellIndex)
	}
	if tree.IsDormant {
		return nil, fmt.Errorf("tree at %v is dormant", cell.Index)
	}
	if tree.Size >= SizeLarge {
		return nil, fmt.Errorf("tree at %v is too large to grow", cell.Index)
	}
	cost := growthCost(tree.IsMine, state, &tree)
	sun := state.MySun
	if sun < cost {
		return nil, fmt.Errorf("player can't afford to grow tree at %v: they have %v but need %v", cell.Index, sun, cost)
	}

	state.MySun = sun - cost
	state.Trees[cell.Index].Size++
	state.Trees[cell.Index].IsDormant = true
	return state, nil
}

func growthCost(isMine bool, state *State, tree *Tree) int8 {
	targetSize := tree.Size + 1
	if targetSize > SizeLarge {
		return CostLifecycleEnd
	} else {
		return cost(isMine, targetSize, state)
	}
}

func cost(isMine bool, size int8, state *State) int8 {
	c := treeBaseCost[size]
	for i := range state.Trees {
		tree := state.Trees[i]
		if tree.Size == size && isMine {
			c++
		}
	}
	return c
}

// ApplyComplete completes a tree
func ApplyComplete(s *State, action *Action) (*State, error) {
	state := s.Clone()

	//coord := board.Coords[action.TargetCellIdx]
	cell := state.Board.Cells[action.TargetCellIdx]
	tree := state.Trees[cell.Index]
	if !tree.Exists {
		return nil, fmt.Errorf("tree at %v not found", cell.Index)
	}
	if !tree.IsMine {
		return nil, fmt.Errorf("tree at %v not owned by player", tree.CellIndex)
	}
	if tree.IsDormant {
		return nil, fmt.Errorf("tree at %v is dormant", cell.Index)
	}
	if tree.Size < SizeLarge {
		return nil, fmt.Errorf("tree at %v is too small to sell", cell.Index)
	}
	cost := growthCost(tree.IsMine, state, &tree)
	sun := state.MySun
	if sun < cost {
		return nil, fmt.Errorf("player can't afford to sell tree at %v: they have %v but need %v", cell.Index, sun, cost)
	}

	state.MySun = sun - cost

	points := state.Nutrients
	if cell.Richness == RichnessMedium {
		points += RichnessBonusMedium
	} else if cell.Richness == RichnessHigh {
		points += RichnessBonusHigh
	}
	state.MyScore += int(points)
	state.Trees[cell.Index].Exists = false
	ClearBit(state.MyTrees, cell.Index)
	ClearBit(state.Larges, cell.Index)
	return state, nil
}

// ApplySeed makes the player cast a seed
func ApplySeed(s *State, action *Action) (*State, error) {
	state := s.Clone()
	targetCoord := state.Board.Coords[action.TargetCellIdx]
	sourceCoord := state.Board.Coords[action.SourceCellIdx]

	targetCell := state.Board.Cells[action.TargetCellIdx]
	sourceCell := state.Board.Cells[action.SourceCellIdx]

	if tree := state.Trees[targetCell.Index]; tree.Exists {
		return nil, fmt.Errorf("target cell %v is not empty", targetCell.Index)
	}
	sourceTree := state.Trees[sourceCell.Index]
	if !sourceTree.Exists {
		return nil, fmt.Errorf("no tree at source %v", sourceCell.Index)
	}
	if sourceTree.Size == SizeSeed {
		return nil, fmt.Errorf("tree at %v is too small to seed", sourceCell.Index)
	}
	if !sourceTree.IsMine {
		return nil, fmt.Errorf("tree %v is not owned by player", sourceCell.Index)
	}
	if sourceTree.IsDormant {
		return nil, fmt.Errorf("tree at %v is dormant", sourceCell.Index)
	}

	distance := sourceCoord.DistanceTo(targetCoord)
	if distance > sourceTree.Size {
		return nil, fmt.Errorf("tree at %v can't seed that far (%v spaces)", sourceCell.Index, distance)
	}
	if targetCell.Richness == RichnessUnusable {
		return nil, fmt.Errorf("target cell %v is unusable", targetCell.Index)
	}

	costOfSeed := cost(true, SizeSeed, state)
	sun := state.MySun
	if sun < costOfSeed {
		return nil, fmt.Errorf("player doesn't have enough energy to seed (%v/%v)", sun, costOfSeed)
	}
	sourceTree.IsDormant = true
	state.MySun -= costOfSeed
	tree := Tree{
		CellIndex: targetCell.Index,
		Size:      SizeSeed,
		IsMine:    true,
		IsDormant: true,
	}
	state.AddTree(tree)
	return state, nil
}
