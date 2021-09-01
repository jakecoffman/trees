package server

import (
	"github.com/jakecoffman/trees/server/game"
	"sync"
)

type GameWrapper struct {
	Code       string
	State      *game.State
	Players    []*Player
	Spectators []*Player
}

func NewGameWrapper(player *Player) *GameWrapper {
	g := &GameWrapper{
		Code:       EncodeToString(6),
		State:      game.New(),
		Players:    []*Player{player},
		Spectators: []*Player{},
	}
	mutex.Lock()
	games[g.Code] = g
	mutex.Unlock()
	return g
}

func (g *GameWrapper) Quit(quitter *Player) {
	// TODO inform players game is over
	_ = quitter

	mutex.Lock()
	for _, p := range g.Players {
		players[p.id].game = nil
	}
	for _, p := range g.Spectators {
		players[p.id].game = nil
	}
	delete(games, g.Code)
	mutex.Unlock()
}

func (g *GameWrapper) Join(player *Player) {
	if len(g.Players) == 2 {
		g.Spectators = append(g.Spectators, player)
	} else {
		g.Players = append(g.Players, player)
	}
}

var mutex = sync.RWMutex{}
var players = map[string]*Player{}
var games = map[string]*GameWrapper{}
