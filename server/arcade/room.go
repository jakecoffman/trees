package arcade

import (
	"fmt"
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

	input chan Cmd
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
				return
			}

			if r.State.Day == maxDays-1 {
				r.sendAll(PlayerMessage{Kind: "msg", Value: "Final turn!"})
			}

			if moves[0].Type == game.Wait && moves[1].Type == game.Wait {
				moves[0] = nil
				moves[1] = nil
			} else {
				if moves[0].Type != game.Wait {
					moves[0] = nil
				}
				if moves[1].Type != game.Wait {
					moves[1] = nil
				}
			}
			// update everyone
			r.sendAllGame()
		}
		select {
		case cmd := <-r.input:
			switch cmd.Kind {
			case CmdJoin:
				r.join(cmd.Player)
				r.sendAllGame()
			case CmdRejoin:
				r.sendAll(PlayerMessage{Kind: "msg", Value: "Player reconnected"})
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
						continue
					}
				} else if cmd.Kind == CmdSeed {
					move = &game.Action{Type: game.Seed, SourceCellIdx: cmd.Src, TargetCellIdx: cmd.Tgt}
					_, err := r.State.ApplySeed(playerIndex, move)
					if err != nil {
						SendMsg(cmd.Player.ws, err.Error())
						continue
					}
				} else if cmd.Kind == CmdGrow {
					move = &game.Action{Type: game.Grow, TargetCellIdx: cmd.Tgt}
					_, err := r.State.ApplyGrow(playerIndex, move)
					if err != nil {
						SendMsg(cmd.Player.ws, err.Error())
						continue
					}
				}

				//log.Println(r.Code, "PLAYER", playerIndex, move.String())
				moves[playerIndex] = move
			}
		}
	}
}

const maxDays = 26

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
	r.sendAll(PlayerMessage{Kind: "msg", Value: "Player connected"})
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
