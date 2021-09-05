package arcade

import (
	"encoding/json"
	"github.com/jakecoffman/trees/server/lib"
	"log"
)

type Player struct {
	id string
	//Name string
	ws   *lib.SafetySocket
	Room *Room `json:"-"`
}

func (p *Player) MarshalJSON() ([]byte, error) {
	var connected bool
	if p.id == "bot" {
		connected = true
	} else {
		connected = p.ws != nil
	}
	return json.Marshal(struct {
		Name      string
		Connected bool
	}{
		Name:      "Bot",
		Connected: connected,
	})
}

func (p *Player) Unlock() {
	if p.ws != nil {
		if err := p.ws.WriteJSON(PlayerMessage{Kind: "unlock"}); err != nil {
			log.Println(err)
			return
		}
	}
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
