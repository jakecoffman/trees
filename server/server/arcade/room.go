package arcade

import (
	"github.com/jakecoffman/trees/server/game"
	"log"
)

//
//const (
//	CmdJoin = iota + 1
//	CmdQuit
//	CmdEnd
//	CmdSeed
//	CmdGrow
//	CmdSell
//)

type Room struct {
	Code       string
	State      *game.State
	Players    []*Player
	Spectators []*Player

	input chan interface{}
}

func NewRoom(player *Player) *Room {
	g := &Room{
		Code:       EncodeToString(6),
		State:      game.New(),
		Players:    []*Player{player},
		Spectators: []*Player{},
		input:      make(chan interface{}),
	}
	Building.CreateGame(g)
	go g.loop()
	return g
}

type PlayerMessage struct {
	Kind  string
	Value string `json:",omitempty"`
	Room  *Room  `json:",omitempty"`
	You   int
}

func (r *Room) sendAllGame() {
	arr := append(r.Players, r.Spectators...)
	for i, p := range arr {
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
	arr := append(r.Players, r.Spectators...)
	for _, p := range arr {
		if p.ws != nil {
			if err := p.ws.WriteJSON(msg); err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func (r *Room) Command(cmd interface{}) {
	r.input <- cmd
}

func (r *Room) Join(joiner *Player) {
	r.input <- CmdJoin{joiner}
}

func (r *Room) Rejoin(player *Player) {
	r.input <- CmdRejoin{player}
}

func (r *Room) Quit(quitter *Player) {
	r.input <- CmdQuit{quitter}
}

func (r *Room) loop() {
	for {
		r.sendAllGame()

		select {
		case cmd := <-r.input:
			switch c := cmd.(type) {
			case CmdJoin:
				r.join(c.Player)
			case CmdQuit:

			}
		}
	}
}

type CmdJoin struct {
	Player *Player
}

type CmdRejoin struct {
	Player *Player
}

type CmdQuit struct {
	Player *Player
}

func (r *Room) join(player *Player) {
	if len(r.Players) == 2 {
		r.Spectators = append(r.Spectators, player)
	} else {
		r.Players = append(r.Players, player)
	}
	r.sendAll(PlayerMessage{Kind: "msg", Value: "Player connected"})
}

func (r *Room) rejoin(player *Player) {
	r.sendAll(PlayerMessage{Kind: "msg", Value: "Player reconnected"})
}

func (r *Room) quit(quitter *Player) {
	for _, p := range r.Players {
		if p == quitter {
			r.sendAll(PlayerMessage{Kind: "msg", Value: "Player has quit"})
			return
		}
	}
	for i, p := range r.Spectators {
		if p.id == quitter.id {
			r.Spectators = append(r.Spectators[i:], r.Spectators[i+1:]...)
			break
		}
	}
	Building.Shut(r)
}
