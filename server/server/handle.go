package server

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Kind  string
	Value string       `json:",omitempty"`
	Game  *GameWrapper `json:",omitempty"`
	You   int          `json:",omitempty"`
}

func Handle(ws *websocket.Conn, r *http.Request) {
	var playerId string
	var player *Player

	// extract cookie or create a new playerId
	{
		cookie, err := r.Cookie("player")
		log.Println(cookie, err, r.Cookies())
		if err == http.ErrNoCookie {
			return
		} else {
			playerId = cookie.Value
		}
	}

	// grab player object or create a new one
	mutex.RLock()
	player = players[playerId]
	mutex.RUnlock()
	if player == nil {
		mutex.Lock()
		players[playerId] = &Player{
			id: playerId,
			ws: ws,
		}
		mutex.Unlock()
	}

	// handle what the player is trying to do
	action := r.URL.Query().Get("action")
	code := r.URL.Query().Get("code")

	switch action {
	case "new":
		if player.game != nil {
			player.game.Quit(player)
			player.game = nil
		}
		player.game = NewGameWrapper(player)
	case "join":
		if player.game != nil {
			player.game.Quit(player)
			player.game = nil
		}
		mutex.RLock()
		g, ok := games[code]
		mutex.RUnlock()
		if !ok {
			// TODO tell player the game is gone or code is wrong
			return
		}
		g.Join(player)
		player.game = g
	default:
		return
	}

	sendAll(player.game.Players)

	for {
		if err := loop(player); err != nil {
			break
		}
	}
}

func sendAll(to map[string]*Player) {
	for _, p := range to {
		msg := Message{
			Kind: "game",
			Game: p.game,
		}
		if err := p.ws.WriteJSON(msg); err != nil {
			return
		}
	}
}

func loop(p *Player) error {
	mt, message, err := p.ws.ReadMessage()
	if err != nil {
		//log.Println("read:", err)
		return err
	}
	//log.Printf("recv: %s", message)
	err = p.ws.WriteMessage(mt, message)
	if err != nil {
		//log.Println("write:", err)
		return err
	}
	return nil
}
