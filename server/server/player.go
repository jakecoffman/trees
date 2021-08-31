package server

import (
	"github.com/gorilla/websocket"
)

type Player struct {
	id   string
	ws   *websocket.Conn
	game *GameWrapper
}
