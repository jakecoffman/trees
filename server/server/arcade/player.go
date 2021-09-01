package arcade

import (
	"encoding/json"
	"github.com/jakecoffman/trees/server/server/lib"
	"log"
)

type Player struct {
	id   string
	Name string
	You  bool
	ws   *lib.SafetySocket
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

func SendMsg(ws *lib.SafetySocket, text string) {
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
