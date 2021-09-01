package server

import (
	"encoding/json"
	"github.com/gorilla/websocket"
)

type Player struct {
	id   string
	Name string
	ws   *websocket.Conn
	game *GameWrapper
}

func (p *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		Connected bool
	}{
		Name:      p.Name,
		Connected: p.ws != nil,
	})
}
