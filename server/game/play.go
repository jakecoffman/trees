package game

import (
	"fmt"
	"math/rand"
)

func (state *State) InitStartingTrees() {
	startingCoords := state.Board.Edges()

	// remove if richness is null
	{
		var dedupeStartingCoords []Coord
		for _, c := range startingCoords {
			if state.Board.Map[c].Richness != RichnessUnusable {
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

	// this is done in New but done here just in case it gets called again?
	state.Trees = map[int]*Tree{}
	for i := 0; i < startingTreeCount; i++ {
		{
			cell := state.Board.Map[validCoords[2*i]]
			state.Trees[cell.Index] = &Tree{
				CellIndex: cell.Index,
				Size:      SizeSmall,
				Owner:     0,
				IsDormant: false,
			}
		}
		{
			cell := state.Board.Map[validCoords[2*i+1]]
			state.Trees[cell.Index] = &Tree{
				CellIndex: cell.Index,
				Size:      SizeSmall,
				Owner:     1,
				IsDormant: false,
			}
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
func (state *State) RandomizeBoard() {
	// reset all cells to the proper richness
	for coord, cell := range state.Board.Map {
		cell.Neighbors = state.Board.GetNeighborIds(coord)
		state.Board.Cells[cell.Index] = cell
		state.Board.Map[coord] = cell
	}

	wantedEmptyCells := rand.Intn(maxEmptyCells + 1)
	actuallyEmptyCells := 0
	for actuallyEmptyCells < wantedEmptyCells-1 {
		index := rand.Intn(len(state.Board.Map))
		randCoord := state.Board.Coords[index]
		if state.Board.Map[randCoord].Richness != RichnessUnusable {
			cell := state.Board.Map[randCoord]
			cell.Richness = RichnessUnusable
			state.Board.Map[randCoord] = cell
			state.Board.Cells[cell.Index].Richness = RichnessUnusable
			actuallyEmptyCells++
			if randCoord != randCoord.Opposite() {
				cell := state.Board.Map[randCoord]
				cell.Richness = RichnessUnusable
				state.Board.Map[randCoord] = cell
				actuallyEmptyCells++
				state.Board.Cells[cell.Index].Richness = RichnessUnusable
			}
		}
	}
}

// IsSet returns true iff the nth bit of v is set
func IsSet(v uint64, n int) bool {
	return (v & 1 << n) > 0
}

// GatherSun calculates shadows and adds energy to players
func (state *State) GatherSun() {
	state.Shadows = state.Sun.CalculateShadows(state)
	for i := range state.Trees {
		tree := state.Trees[i]
		state.Trees[i].IsDormant = false
		if tree.Owner == 0 {
			if !IsSet(state.Shadows[0][tree.Size], tree.CellIndex) {
				state.Energy[0] += tree.Size
			}
		} else {
			if !IsSet(state.Shadows[0][tree.Size], tree.CellIndex) {
				state.Energy[1] += tree.Size
			}
		}
	}
}

var treeBaseCost = []int{0, 1, 3, 7}

const (
	CostLifecycleEnd    = 4
	RichnessBonusMedium = 2
	RichnessBonusHigh   = 4
)

// ApplyGrow makes player grow a tree
func (state *State) ApplyGrow(player int, action Action) error {
	cell := state.Board.Cells[action.TargetCellIdx]
	tree, ok := state.Trees[cell.Index]
	if !ok {
		return fmt.Errorf("tree not found at %v", cell.Index)
	}
	if player != tree.Owner {
		return fmt.Errorf("tree at %v not owned by player %v", tree.CellIndex, player)
	}
	if tree.IsDormant {
		return fmt.Errorf("tree at %v is dormant", cell.Index)
	}
	if tree.Size >= SizeLarge {
		return fmt.Errorf("tree at %v is too large to grow", cell.Index)
	}
	cost := growthCost(player, state, tree)
	sun := state.Energy[player]
	if sun < cost {
		return fmt.Errorf("player %v can't afford to grow tree at %v: they have %v but need %v", player, cell.Index, sun, cost)
	}

	state.Energy[player] = sun - cost
	state.Trees[cell.Index].Size++
	state.Trees[cell.Index].IsDormant = true
	return nil
}

func growthCost(player int, state *State, tree *Tree) int {
	targetSize := tree.Size + 1
	if targetSize > SizeLarge {
		return CostLifecycleEnd
	} else {
		return cost(player, targetSize, state)
	}
}

func cost(player int, size int, state *State) int {
	c := treeBaseCost[size]
	for i := range state.Trees {
		tree := state.Trees[i]
		if tree.Size == size && player == tree.Owner {
			c++
		}
	}
	return c
}

// ApplyComplete completes a tree
func ApplyComplete(player int, state *State, action Action) error {
	//coord := board.Coords[action.TargetCellIdx]
	cell := state.Board.Cells[action.TargetCellIdx]
	tree, ok := state.Trees[cell.Index]
	if !ok {
		return fmt.Errorf("tree at %v not found", cell.Index)
	}
	if player != tree.Owner {
		return fmt.Errorf("tree at %v not owned by player %v", tree.CellIndex, player)
	}
	if tree.IsDormant {
		return fmt.Errorf("tree at %v is dormant", cell.Index)
	}
	if tree.Size < SizeLarge {
		return fmt.Errorf("tree at %v is too small to sell", cell.Index)
	}
	cost := growthCost(player, state, tree)
	sun := state.Energy[player]
	if sun < cost {
		return fmt.Errorf("player %v can't afford to sell tree at %v: they have %v but need %v", player, cell.Index, sun, cost)
	}

	state.Energy[player] = sun - cost

	points := state.Nutrients
	if cell.Richness == RichnessMedium {
		points += RichnessBonusMedium
	} else if cell.Richness == RichnessHigh {
		points += RichnessBonusHigh
	}
	state.Score[player] += points
	delete(state.Trees, cell.Index)
	return nil
}

// ApplySeed makes the player cast a seed
func ApplySeed(player int, state *State, action Action) error {
	targetCoord := state.Board.Coords[action.TargetCellIdx]
	sourceCoord := state.Board.Coords[action.SourceCellIdx]

	targetCell := state.Board.Cells[action.TargetCellIdx]
	sourceCell := state.Board.Cells[action.SourceCellIdx]

	if _, ok := state.Trees[targetCell.Index]; ok {
		return fmt.Errorf("target cell %v is not empty", targetCell.Index)
	}
	sourceTree := state.Trees[sourceCell.Index]
	if sourceTree == nil {
		return fmt.Errorf("no tree at source %v", sourceCell.Index)
	}
	if sourceTree.Size == SizeSeed {
		return fmt.Errorf("tree at %v is too small to seed", sourceCell.Index)
	}
	if player != sourceTree.Owner {
		return fmt.Errorf("tree %v is not owned by player %v", sourceCell.Index, player)
	}
	if sourceTree.IsDormant {
		return fmt.Errorf("tree at %v is dormant", sourceCell.Index)
	}

	distance := sourceCoord.DistanceTo(targetCoord)
	if distance > sourceTree.Size {
		return fmt.Errorf("tree at %v can't seed that far (%v spaces)", sourceCell.Index, distance)
	}
	if targetCell.Richness == RichnessUnusable {
		return fmt.Errorf("target cell %v is unusable", targetCell.Index)
	}

	costOfSeed := cost(player, SizeSeed, state)
	sun := state.Energy[player]
	if sun < costOfSeed {
		return fmt.Errorf("player %v doesn't have enough energy to seed (%v/%v)", player, sun, costOfSeed)
	}
	sourceTree.IsDormant = true
	state.Energy[player] -= costOfSeed
	tree := Tree{
		CellIndex: targetCell.Index,
		Size:      SizeSeed,
		Owner:     player,
		IsDormant: true,
	}
	state.Trees[tree.CellIndex] = &tree
	return nil
}
