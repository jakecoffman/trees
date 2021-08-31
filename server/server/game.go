package server

import (
	"github.com/jakecoffman/trees/server/game"
	"sync"
)

type GameWrapper struct {
	Code       string
	State      *game.State
	Players    map[string]*Player
	Spectators map[string]*Player
}

func NewGameWrapper(player *Player) *GameWrapper {
	g := &GameWrapper{
		Code:       EncodeToString(6),
		State:      game.New(),
		Players:    map[string]*Player{player.id: player},
		Spectators: map[string]*Player{},
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
		g.Spectators[player.id] = player
	} else {
		g.Players[player.id] = player
	}
}

var mutex = sync.RWMutex{}
var players = map[string]*Player{}
var games = map[string]*GameWrapper{}
