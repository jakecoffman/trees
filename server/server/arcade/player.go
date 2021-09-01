package arcade

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type Player struct {
	id   string
	Name string
	You  bool
	ws   *websocket.Conn
	Room *Room `json:"-"`
}

func (p *Player) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Name      string
		You       bool
		Connected bool
	}{
		Name:      p.Name,
		You:       p.You,
		Connected: p.ws != nil,
	})
}

func SendMsg(ws *websocket.Conn, text string) {
	msg := PlayerMessage{
		Kind:  "msg",
		Value: text,
	}
	if ws != nil {
		if err := ws.WriteJSON(msg); err != nil {
			log.Println(err)
			return
		}
	} else {
		log.Println("SKIPPED")
	}
}
