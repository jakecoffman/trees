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
		if err != nil {
			log.Println("Player failed to connect:", err.Error())
			return
		} else {
			playerId = cookie.Value
			log.Println("Reconnect", playerId)
		}
	}

	// grab player object or create a new one
	mutex.RLock()
	player = players[playerId]
	mutex.RUnlock()
	if player == nil {
		mutex.Lock()
		player = &Player{
			id: playerId,
			ws: ws,
		}
		players[playerId] = player
		mutex.Unlock()
	} else {
		player.ws = ws
	}
	defer func() {
		mutex.Lock()
		player.ws = nil
		mutex.Unlock()
	}()

	// handle what the player is trying to do
	log.Println("QUERY", r.URL.Query())
	action := r.URL.Query().Get("action")
	code := r.URL.Query().Get("code")

	switch action {
	case "new":
		log.Println("Starting new game")
		if player.game != nil {
			player.game.Quit(player)
			player.game = nil
		}
		player.game = NewGameWrapper(player)
	case "join":
		if player.game != nil && player.game.Code != code {
			player.game.Quit(player)
			player.game = nil
		}
		if player.game == nil {
			mutex.RLock()
			g, ok := games[code]
			mutex.RUnlock()
			if !ok {
				log.Printf("Code is wrong, or game is done: %v - %v\n", code, games)
				// TODO tell player the game is gone or code is wrong
				return
			}
			g.Join(player)
			player.game = g
		}
	default:
		log.Println("Invalid/missing action:", action)
		return
	}

	log.Println("Sending")
	sendAll(player.game.Players...)
	log.Println("Sent")

	for {
		if err := loop(player); err != nil {
			break
		}
	}
}

func sendAll(to ...*Player) {
	for _, p := range to {
		msg := Message{
			Kind: "game",
			Game: p.game,
		}
		if p.ws != nil {
			if err := p.ws.WriteJSON(msg); err != nil {
				log.Println(err)
				return
			}
		} else {
			log.Println("SKIPPED")
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
