package bot

import (
	"log"
	"math/rand"
	"runtime/debug"
	"testing"
	"time"
)

func TestOnePlayer(t *testing.T) {
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

	var wins [2]int

	// simulate games
	for i := 0; i < 100; i++ {
		state := &State{
			Board:     NewBoard(),
			Nutrients: 20,
			Trees:     make([]Tree, 37),
			Num:       make([]int8, SizeLarge+1),
		}
		//RandomizeBoard(state)
		InitStartingTrees(state)
		state.Shadows = state.Sun.CalculateShadows(state.Board, state.Trees)
		state = GatherSun(state)
		var err error

		// sanity
		//log.Println(state.MySun)
		//for _, t := range state.Trees {
		//	if t.Exists {
		//		log.Println(t)
		//	}
		//}

		var myLastAction, oppLastAction Action
		for state.Day < 23 {
			log.Println(state.GeneratePossibleMoves())

			var myAction Action
			if myLastAction.Type == Wait && oppLastAction.Type != Wait {
				myAction = myLastAction
			} else {
				myAction = nextAction(state, mySettings)
			}

			log.Println(myAction)

			myLastAction = myAction

			if myAction.Type == Wait {
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
			var nutrientsLost int8
			if myAction.Type == Complete {
				nutrientsLost++
				state, err = ApplyComplete(state, &myAction)
				if err != nil {
					panic(err)
				}
			}
			state.Nutrients = max(0, state.Nutrients-nutrientsLost)
			if myAction.Type == Grow {
				state, err = ApplyGrow(state, &myAction)
				if err != nil {
					panic(err)
				}
			}
		}

		if state.MyScore > bestScore {
			bestScore = state.MyScore
			//log.Println("HIGH SCORE", bestScore, mySettings)
		}

		if state.MyScore == state.OpponentScore {
			log.Println(i, "SCORE", state.MyScore, "-", state.OpponentScore)
		}

		t.Logf("Wins: %v - %v", wins[0], wins[1])
	}
	t.Log("Best:", bestScore, settings)
}
