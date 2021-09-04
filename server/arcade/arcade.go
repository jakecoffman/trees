package arcade

import (
	"github.com/jakecoffman/trees/server/lib"
	"log"
	"sync"
)

type bldg struct {
	mutex   sync.RWMutex
	players map[string]*Player
	games   map[string]*Room
}

var Building = bldg{
	players: map[string]*Player{},
	games:   map[string]*Room{},
}

// Enter the player enters the Arcade
func (b *bldg) Enter(playerId string, ws *lib.SafetySocket) *Player {
	b.mutex.RLock()
	player := b.players[playerId]
	b.mutex.RUnlock()
	if player == nil {
		b.mutex.Lock()
		log.Println(playerId, "New Player Object")
		player = &Player{
			id: playerId,
			ws: ws,
		}
		b.players[playerId] = player
		b.mutex.Unlock()
	} else {
		player.ws = ws
	}
	return player
}

// Disconnect the player left
func (b *bldg) Disconnect(player *Player) {
	log.Println(player.id, "Disconnected")
	player.ws = nil
	if player.Room != nil {
		// tell the others what has happened
		player.Room.sendAllGame()
		player.Room.sendAll(PlayerMessage{Kind: "msg", Value: "Player disconnected"})
	}
}

func (b *bldg) CreateGame(g *Room) {
	b.mutex.Lock()
	b.games[g.Code] = g
	b.mutex.Unlock()
}

// Shut closes the room so players don't enter it anymore
func (b *bldg) Shut(r *Room) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	delete(b.games, r.Code)
}

func (b *bldg) FindRoom(code string) (*Room, bool) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	g, ok := b.games[code]
	return g, ok
}
