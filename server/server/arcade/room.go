package arcade

import (
	"github.com/jakecoffman/trees/server/game"
	"log"
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
}

func (r *Room) Rejoin(player *Player) {
	r.input <- Cmd{Kind: CmdRejoin, Player: player}
}

func (r *Room) Quit(quitter *Player) {
	r.input <- Cmd{Kind: CmdQuit, Player: quitter}
}

func (r *Room) CastSeed(player *Player, source, target int) {
	r.input <- Cmd{Kind: CmdSeed, Player: player, Src: source, Tgt: target}
}

func (r *Room) SellTree(player *Player, index int) {
	r.input <- Cmd{Kind: CmdSell, Player: player, Src: index}
}

func (r *Room) GrowTree(player *Player, index int) {
	r.input <- Cmd{Kind: CmdGrow, Player: player, Src: index}
}

func (r *Room) EndTurn(player *Player) {
	r.input <- Cmd{Kind: CmdEnd, Player: player}
}

func (r *Room) loop() {
	var moves [2]*game.Action

	for {
		if moves[0] != nil && moves[1] != nil {
			log.Println("BOTH PLAYERS ENTERED A MOVE")
			// both players have entered a move, execute them "simultaneously"

			moves[0] = nil
			moves[1] = nil
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
			default:
				var move *game.Action

				if cmd.Kind == CmdEnd {
					move = &game.Action{Type: game.Wait}
				} else if cmd.Kind == CmdSell {
					move = &game.Action{Type: game.Complete, SourceCellIdx: cmd.Src}
				} else if cmd.Kind == CmdSeed {
					move = &game.Action{Type: game.Seed, SourceCellIdx: cmd.Src, TargetCellIdx: cmd.Tgt}
				} else if cmd.Kind == CmdGrow {
					move = &game.Action{Type: game.Grow, SourceCellIdx: cmd.Src}
				}

				if cmd.Player == r.Players[0] {
					log.Println(r.Code, "PLAYER 0", move.String())
					moves[0] = move
				} else {
					log.Println(r.Code, "PLAYER 1", move.String())
					moves[1] = move
				}
			}
		}
	}
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
	for _, p := range r.Players {
		if p == quitter {
			r.sendAll(PlayerMessage{Kind: "msg", Value: "Player has quit"})
			return
		}
	}
	Building.Shut(r)
}
