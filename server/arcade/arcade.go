package arcade

import (
	"github.com/jakecoffman/trees/server/lib"
	"log"
	"sync"
)

var ActiveWsConnections int64

type bldg struct {
	mutex   sync.RWMutex
	players map[string]*Player
	games   map[string]*Room
}

var Building = bldg{
	players: map[string]*Player{},
	games:   map[string]*Room{},
}

func (b *bldg) Counts() (int, int) {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return len(b.players), len(b.games)
}

func (b *bldg) GetPlayer(playerId string) *Player {
	b.mutex.RLock()
	defer b.mutex.RUnlock()
	return b.players[playerId]
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
		//player.Room.sendAll(PlayerMessage{Kind: "msg", Value: "Player disconnected"})
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
	if _, ok := b.games[r.Code]; ok {
		delete(b.games, r.Code)
	}
}

func (b *bldg) GetRoom(code string) *Room {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	g := b.games[code]
	return g
}
