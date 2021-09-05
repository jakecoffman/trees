package arcade

import (
	"fmt"
	"github.com/jakecoffman/trees/server/bot"
	"github.com/jakecoffman/trees/server/game"
	"log"
	"runtime/debug"
)

const (
	CmdJoin = iota + 1
	CmdRejoin
	CmdQuit
	CmdEnd
	CmdSeed
	CmdGrow
	CmdSell
)

type Room struct {
	Code    string
	State   *game.State
	Players []*Player

	input                chan Cmd
	finalDayNotification bool
	playingBot           bool
}

func NewRoom() *Room {
	g := &Room{
		Code:  EncodeToString(6),
		State: game.New(),
		input: make(chan Cmd),
	}
	Building.CreateGame(g)
	log.Println("New room created", g.Code)
	go g.loop()
	return g
}

type PlayerMessage struct {
	Kind   string
	Value  string `json:",omitempty"`
	Source int    `json:",omitempty"`
	Target int    `json:",omitempty"`
	Room   *Room  `json:",omitempty"`
	You    int
}

func (r *Room) sendAllGame() {
	for i, p := range r.Players {
		if p.ws != nil {
			if err := p.ws.WriteJSON(PlayerMessage{
				Kind: "room",
				Room: r,
				You:  i,
			}); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (r *Room) sendAll(msg PlayerMessage) {
	for _, p := range r.Players {
		if p.ws != nil {
			if err := p.ws.WriteJSON(msg); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (r *Room) Join(joiner *Player) {
	r.input <- Cmd{Kind: CmdJoin, Player: joiner}
	joiner.Room = r
}

func (r *Room) Rejoin(player *Player) {
	r.input <- Cmd{Kind: CmdRejoin, Player: player}
	player.Room = r
}

func (r *Room) Quit(quitter *Player) {
	r.input <- Cmd{Kind: CmdQuit, Player: quitter}
	quitter.Room = nil
}

func (r *Room) CastSeed(player *Player, source, target int) {
	r.input <- Cmd{Kind: CmdSeed, Player: player, Src: source, Tgt: target}
}

func (r *Room) SellTree(player *Player, index int) {
	r.input <- Cmd{Kind: CmdSell, Player: player, Tgt: index}
}

func (r *Room) GrowTree(player *Player, index int) {
	r.input <- Cmd{Kind: CmdGrow, Player: player, Tgt: index}
}

func (r *Room) EndTurn(player *Player) {
	r.input <- Cmd{Kind: CmdEnd, Player: player}
}

func (r *Room) loop() {
	var moves [2]*game.Action
	defer func() {
		if rc := recover(); rc != nil {
			fmt.Println("Room", r.Code, "has crashed:", rc, string(debug.Stack()))
		} else {
			log.Println("Room", r.Code, "has closed")
		}
	}()

	for {
		if moves[0] != nil && moves[1] != nil {
			// both players have entered a move, execute them "simultaneously"
			var isGameOver bool
			if moves, isGameOver = r.execute(moves); isGameOver {
				return
			}
		}
		if r.playingBot && moves[1] == nil {
			move := doBotMove(r.State)
			moves[1] = move
			if moves[0] != nil {
				continue
			}
		}
		select {
		case cmd := <-r.input:
			switch cmd.Kind {
			case CmdJoin:
				r.join(cmd.Player)
				r.sendAllGame()
			case CmdRejoin:
				//r.sendAll(PlayerMessage{Kind: "msg", Value: "Player reconnected"})
				r.sendAllGame()
			case CmdQuit:
				r.quit(cmd.Player)
				r.sendAllGame()
				return // terminates the game loop
			default:
				var move *game.Action
				var playerIndex int

				if cmd.Player == r.Players[0] {
					playerIndex = 0
				} else {
					playerIndex = 1
				}

				if cmd.Kind == CmdEnd {
					move = &game.Action{Type: game.Wait}
				} else if cmd.Kind == CmdSell {
					move = &game.Action{Type: game.Complete, TargetCellIdx: cmd.Tgt}
					_, err := r.State.ApplyComplete(playerIndex, move)
					if err != nil {
						SendMsg(cmd.Player.ws, err.Error())
						cmd.Player.Unlock()
						continue
					}
				} else if cmd.Kind == CmdSeed {
					move = &game.Action{Type: game.Seed, SourceCellIdx: cmd.Src, TargetCellIdx: cmd.Tgt}
					_, err := r.State.ApplySeed(playerIndex, move)
					if err != nil {
						SendMsg(cmd.Player.ws, err.Error())
						cmd.Player.Unlock()
						continue
					}
				} else if cmd.Kind == CmdGrow {
					move = &game.Action{Type: game.Grow, TargetCellIdx: cmd.Tgt}
					_, err := r.State.ApplyGrow(playerIndex, move)
					if err != nil {
						SendMsg(cmd.Player.ws, err.Error())
						cmd.Player.Unlock()
						continue
					}
				}

				//log.Println(r.Code, "PLAYER", playerIndex, move.String())
				moves[playerIndex] = move
			}
		}
	}
}

func doBotMove(state *game.State) *game.Action {
	botState := &bot.State{
		Board:         bot.NewBoard(),
		Sun:           bot.Sun{Orientation: int8(state.Sun.Orientation)},
		Day:           int8(state.Day),
		Nutrients:     int8(state.Nutrients),
		MySun:         int8(state.Energy[1]),
		OpponentSun:   int8(state.Energy[0]),
		MyScore:       state.Score[1],
		OpponentScore: state.Score[0],
		Trees:         make([]bot.Tree, 37),
		Num:           make([]int8, bot.SizeLarge+1),
	}
	for _, v := range state.Trees {
		botState.AddTree(bot.Tree{
			CellIndex: int8(v.CellIndex),
			Size:      int8(v.Size),
			IsMine:    v.Owner == 1,
			IsDormant: v.IsDormant,
			Exists:    true,
		})
	}
	// copy over the unusables
	for i, v := range state.Board.Cells {
		if v.Richness == game.RichnessUnusable {
			botState.Board.Cells[i].Richness = game.RichnessUnusable
		}
	}
	botState.Shadows = botState.Sun.CalculateShadows(botState.Board, botState.Trees)
	nextMove, _, _ := bot.Chokudai(botState, &bot.Settings{
		MinimalEconomy: 16,
		MaximumEconomy: 25,
	})
	move := &game.Action{
		Type:          int(nextMove.Type),
		TargetCellIdx: int(nextMove.TargetCellIdx),
		SourceCellIdx: int(nextMove.SourceCellIdx),
	}
	return move
}

func (r *Room) execute(moves [2]*game.Action) ([2]*game.Action, bool) {
	r.applyMoves(moves)

	if r.State.Day == maxDays {
		r.State.Score[0] += r.State.Energy[0] / 3
		r.State.Score[1] += r.State.Energy[1] / 3
		r.sendAllGame()
		var msg string
		if r.State.Score[0] == r.State.Score[1] {
			msg = "Tie!"
		} else if r.State.Score[0] > r.State.Score[1] {
			msg = "Orange wins!"
		} else {
			msg = "Blue wins!"
		}
		r.sendAll(PlayerMessage{
			Kind:  "msg",
			Value: msg,
		})
		Building.Shut(r)
		r.Players[0].Room = nil
		r.Players[1].Room = nil
		return [2]*game.Action{}, true
	}

	if r.State.Day == maxDays-1 && !r.finalDayNotification {
		r.sendAll(PlayerMessage{Kind: "msg", Value: "Final turn!"})
		r.finalDayNotification = true
	}

	// if they are both wait, unlock them since the day has progressed
	if moves[0].Type == game.Wait && moves[1].Type == game.Wait {
		moves[0] = nil
		moves[1] = nil
		r.sendAll(PlayerMessage{Kind: "unlock"})
		log.Println("BOTH WERE WAIT")
	} else {
		// otherwise unlock them if they did not wait, so the waiter keeps waiting
		if moves[0].Type != game.Wait {
			moves[0] = nil
			r.Players[0].Unlock()
			log.Println("P1 WAIT")
		}
		if moves[1].Type != game.Wait {
			moves[1] = nil
			r.Players[1].Unlock()
			log.Println("P2 WAIT")
		}
	}
	// update everyone
	r.sendAllGame()
	return moves, false
}

const maxDays = 24

func (r *Room) applyMoves(moves [2]*game.Action) {
	state := r.State

	if moves[0].Type == game.Seed && moves[1].Type == game.Seed {
		if moves[0].TargetCellIdx == moves[1].TargetCellIdx {
			// seed is not planted, refund sun, make source trees dormant
			state.Trees[moves[0].SourceCellIdx].IsDormant = true
			state.Trees[moves[1].SourceCellIdx].IsDormant = true
			r.sendAll(PlayerMessage{Kind: "msg", Value: "Both players seeded the same hex, seed was removed."})
			// TODO refund sun?
			return
		}
	}

	if moves[0].Type == game.Wait && moves[1].Type == game.Wait {
		state.Day++
		state.Sun = state.Sun.Move()
		state = state.GatherSun()
	}
	if moves[0].Type == game.Seed {
		state = sanity(state.ApplySeed(0, moves[0]))
	}
	if moves[1].Type == game.Seed {
		state = sanity(state.ApplySeed(1, moves[1]))
	}
	var nutrientsLost int
	if moves[0].Type == game.Complete {
		nutrientsLost++
		state = sanity(state.ApplyComplete(0, moves[0]))
	}
	if moves[1].Type == game.Complete {
		nutrientsLost++
		state = sanity(state.ApplyComplete(1, moves[1]))
	}
	state.Nutrients = max(0, state.Nutrients-nutrientsLost)
	if moves[0].Type == game.Grow {
		state = sanity(state.ApplyGrow(0, moves[0]))
	}
	if moves[1].Type == game.Grow {
		state = sanity(state.ApplyGrow(1, moves[1]))
	}
	state.Shadows = state.Sun.CalculateShadows(state)
	r.State = state
}

// there shouldn't be errors, but if there is lets panic so they can be fixed
func sanity(state *game.State, err error) *game.State {
	if err != nil {
		panic(err)
	}
	if state == nil {
		panic("Unexpected nil")
	}
	return state
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Cmd struct {
	Kind   int
	Player *Player
	Src    int
	Tgt    int
}

func (r *Room) join(player *Player) {
	if len(r.Players) == 2 {
		log.Println(r.Code, "Room is full, someone tried to join")
		return
	}
	if len(r.Players) == 1 && player == r.Players[0] {
		log.Println(r.Code, "Room double join?")
		return
	}
	r.Players = append(r.Players, player)
	//r.sendAll(PlayerMessage{Kind: "msg", Value: "Player connected"})
}

func (r *Room) quit(quitter *Player) {
	for i, p := range r.Players {
		if p == quitter {
			r.sendAll(PlayerMessage{Kind: "msg", Value: fmt.Sprintf("Player %v has quit", i+1)})
		}
		p.Room = nil
	}
	Building.Shut(r)
}

func (r *Room) UseBot() {
	if len(r.Players) == 1 {
		log.Println("bot time")
		r.playingBot = true
		r.Join(&Player{
			id:   "bot",
			ws:   nil,
			Room: r,
		})
	}
}
