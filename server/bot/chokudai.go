package bot

import (
	"container/heap"
	"log"
	"strings"
	"syscall"
)

const chokudaiMaxTurns = 10
const chokudaiWidth = 100
const msLimit = 100

// Chokudai search: returns best action and score of the best end state
func Chokudai(first *State, s *Settings) (Action, *State, []*State) {
	var queues = make([]PriorityQueue, chokudaiMaxTurns+1)
	queue := &queues[0]
	heap.Push(queue, &Item{State: first})
	width := chokudaiWidth

	processed := 1

	var realMaxDepth int
	considered := 0
	var tv syscall.Timeval
	_ = syscall.Gettimeofday(&tv)
	start := int64(tv.Sec)*1e3 + int64(tv.Usec)/1e3
	for processed > 0 {
		_ = syscall.Gettimeofday(&tv)
		if (int64(tv.Sec)*1e3+int64(tv.Usec)/1e3)-start > msLimit {
			break
		}
		for t := 0; t < chokudaiMaxTurns; t++ {
			for i := 0; i < width; i++ {
				if queues[t].Len() == 0 {
					break
				}
				current := heap.Pop(&queues[t]).(*Item)
				processed--

				// we've reached the final day
				if current.Day == 24 {
					heap.Push(&queues[chokudaiMaxTurns], current)
					break
				}

				moves := current.State.GeneratePossibleMoves()

				//if bits.OnesCount64(current.MyTrees & current.Seeds & ^current.DormantTrees) >= 1 && current.MySun >= 1 {
				//	for i := 0; i < len(current.Trees); i++ {
				//		t := current.Trees[i]
				//		if t.Exists && t.IsMine && t.Size == SizeSeed {
				//			log.Println(t.IsDormant, current.GeneratedByAction)
				//		}
				//	}
				//	d, _ := json.Marshal(current.State)
				//	log.Println(current.MySun, current.Num, moves, string(d))
				//}

				// prune
				for z := 0; z < len(moves); {
					nextAction := moves[z]
					// Careful, these pruning things break things
					if current.Num[SizeSeed] == 1 && nextAction.Type == Seed {
						// don't seed when you already have a seed
						moves = append(moves[:z], moves[z+1:]...)
					} else if current.Day < 11 && nextAction.Type == Complete {
						// don't sell before day 13
						moves = append(moves[:z], moves[z+1:]...)
					} else if nextAction.Type == Seed && current.Trees[nextAction.SourceCellIdx].Size == SizeSmall {
						// don't seed from small trees
						moves = append(moves[:z], moves[z+1:]...)
					} else if nextAction.Type == Seed && 1 == abs(current.Board.Coords[nextAction.TargetCellIdx].DistanceTo(current.Board.Coords[nextAction.SourceCellIdx])) {
						// don't seed 1 space away
						moves = append(moves[:z], moves[z+1:]...)
					} else {
						z++
					}
				}

				var movesAdded int
				for z := 0; z < len(moves); z++ {
					nextAction := moves[z]
					// final prune: don't wait if you can avoid it
					if current.State == first && nextAction.Type == Wait && len(moves) > 1 {
						continue
					}

					nextState := applyMove(nextAction, current.State)

					movesAdded++
					considered++
					nextState.GeneratedByAction = nextAction
					nextState.CameFrom = current.State
					p := nextState.Priority(s)
					if p == -1 {
						// don't seed next to another one of your own trees
						continue
					}
					heap.Push(&queues[t+1], &Item{State: nextState, Priority: p})
					if t+1 < chokudaiMaxTurns {
						processed++
					}

					if t+1 > realMaxDepth {
						realMaxDepth = t + 1
					}
				}
				if movesAdded == 0 {
					log.Println("WARNING: dead ended at state", current.State)
				}
			}
		}
		//width++
	}

	//if processed == 0 {
	//	log.Println("Processed all moves")
	//} else {
	//	log.Println("Unprocessed:", processed)
	//}

	//log.Println("TIME", time.Now().Sub(start))
	//log.Println("BEST PRIORITY", best.Priority())
	//log.Println("BEST SCORE", best.MyScore)
	//log.Println("DEPTH", realMaxDepth)
	//log.Println("WIDEST", width)
	//log.Println("MOVES CONSIDERED", considered)

	if queues[chokudaiMaxTurns].Len() == 0 {
		log.Println("Couldn't find a move worth doing, sad!")
		return Action{Type: Wait}, nil, nil
	}
	best := heap.Pop(&queues[chokudaiMaxTurns]).(*Item)

	//log.Printf("Best state is (%v): %v\n", best.Priority, best)
	//bestScore := best.Priority
	path := patherize(first, best.State)
	//printPath(patherize(first, best.State))
	//for queues[chokudaiMaxTurns].Len() > 0 && path[0].GeneratedByAction.Type == Wait {
	//	nextItem := heap.Pop(&queues[chokudaiMaxTurns]).(*Item)
	//	if nextItem.Priority < bestScore {
	//		break
	//	}
	//	best = nextItem
	//	path = patherize(first, nextItem.State)
	//}
	//printPath(path)
	return path[0].GeneratedByAction, best.State, path
}

func patherize(first *State, best *State) []*State {
	current := best
	var path []*State
	for current != first {
		path = append(path, current)
		current = current.CameFrom
	}
	if len(path) == 0 {
		log.Fatal("That's not right: ", best.MyScore)
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

func scorePath(first *State, best *State) int {
	current := best
	var path []*State
	for current != first {
		path = append(path, current)
		current = current.CameFrom
	}
	if len(path) == 0 {
		log.Fatal("That's not right: ", best.MyScore)
	}
	var score int
	// remember, it's backwards
	if path[len(path)-1].GeneratedByAction.Type == Wait {
		score -= 10
	}
	return score
}

func printPath(path []*State) {
	var b strings.Builder
	for _, p := range path {
		b.WriteString(p.GeneratedByAction.String())
		b.WriteString(", ")
	}
	log.Println(b.String())
}

func applyMove(nextAction Action, current *State) *State {
	var nextState *State
	switch nextAction.Type {
	case Wait:
		nextState = current.GatherSun()
		//log.Println("WAIT", t, current.State.MyScore, "->", nextState.MyScore)
	case Seed:
		nextState = current.DoSeed(nextAction)
		//log.Println("SEED", t, current.State.MyScore, "->", nextState.MyScore)
	case Complete:
		nextState = current.DoComplete(nextAction)
		//log.Println("COMP", t, current.State.MyScore, "->", nextState.MyScore)
	case Grow:
		nextState = current.DoGrow(nextAction)
		//log.Println("SEED", t, current.State.MyScore, "->", nextState.MyScore)
	default:
		log.Panic("Invalid state: ", nextAction.Type)
	}
	nextState.GeneratedByAction = nextAction
	return nextState
}

func Patherizer(state, first *State) {
	current := state
	var path []*State
	for current != first {
		path = append(path, current)
		current = current.CameFrom
	}
	// reverse
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	for _, s := range path {
		log.Println(s.Priority(&settings), s.GeneratedByAction.String())
	}
}
