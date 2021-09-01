package server

import (
	"fmt"
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
			sendMsg(ws, "Failed to connect, missing player cookie")
			return
		} else {
			playerId = cookie.Value
			log.Println(playerId, "Connected")
		}
	}

	// grab player object or create a new one
	mutex.RLock()
	player = players[playerId]
	mutex.RUnlock()
	if player == nil {
		mutex.Lock()
		log.Println(playerId, "New Player Object")
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
		player.ws = nil
		if player.game != nil {
			sendAllGame(player.game)
			sendAll(player.game, Message{Kind: "msg", Value: "Player disconnected"})
		}
	}()

	// handle what the player is trying to do
	action := r.URL.Query().Get("action")
	code := r.URL.Query().Get("code")

	switch action {
	case "new":
		str := fmt.Sprint(playerId)
		if player.game != nil {
			str += " leaving"
			player.game.Quit(player)
			player.game = nil
		}
		str += " new"
		log.Println(str)
		player.game = NewGameWrapper(player)
	case "join":
		str := fmt.Sprint(playerId)
		if player.game != nil && player.game.Code != code {
			str += " leaving"
			player.game.Quit(player)
			player.game = nil
		}
		if player.game == nil {
			mutex.RLock()
			str += " fresh"
			g, ok := games[code]
			mutex.RUnlock()
			if !ok {
				sendMsg(player.ws, fmt.Sprintf("Code is wrong, or game is gone: %v", code))
				return
			}
			g.Join(player)
			player.game = g
		}
		str += " joined"
		log.Println(str)
	default:
		log.Println("Invalid/missing action:", action)
		sendMsg(ws, "Invalid action: "+action)
		return
	}

	sendAllGame(player.game)
	sendAll(player.game, Message{Kind: "msg", Value: "Player connected"})

	for {
		if err := loop(player); err != nil {
			break
		}
	}
}

func sendAllGame(game *GameWrapper) {
	sendAll(game, Message{
		Kind: "game",
		Game: game,
	})
}

func sendAll(game *GameWrapper, msg Message) {
	arr := append(game.Players, game.Spectators...)
	for _, p := range arr {
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

func sendMsg(ws *websocket.Conn, text string) {
	msg := Message{
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
