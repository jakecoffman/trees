package server

import (
	"fmt"
	"github.com/jakecoffman/trees/server/server/arcade"
	"github.com/jakecoffman/trees/server/server/lib"
	"log"
	"net/http"
	"runtime/debug"
)

func Handle(ws *lib.SafetySocket, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Handler crashed:", r, string(debug.Stack()))
		} else {
			log.Println("Player disconnected")
		}
	}()
	var playerId string

	// extract cookie or create a new playerId
	{
		cookie, err := r.Cookie("player")
		if err != nil {
			log.Println("Player failed to connect:", err.Error())
			arcade.SendMsg(ws, "Failed to connect, missing player cookie")
			return
		} else {
			playerId = cookie.Value
			log.Println(playerId, "Connected")
		}
	}

	player := arcade.Building.Enter(playerId, ws)
	defer arcade.Building.Disconnect(player)

	// handle what the player is trying to do
	action := r.URL.Query().Get("action")
	code := r.URL.Query().Get("code")

	var room *arcade.Room

	switch action {
	case "new":
		str := fmt.Sprint(playerId)
		if player.Room != nil {
			// player is already in a room, leave it
			str += " leaving"
			player.Room.Quit(player)
		}
		str += " new"
		log.Println(str)
		room = arcade.NewRoom()
		room.Join(player)
	case "join":
		str := fmt.Sprint(playerId)
		if player.Room != nil && player.Room.Code != code {
			str += " leaving"
			player.Room.Quit(player)
		}
		if player.Room == nil {
			var ok bool
			room, ok = arcade.Building.FindRoom(code)
			str += " fresh join"
			if !ok {
				arcade.SendMsg(ws, fmt.Sprintf("Code is wrong, or room is gone: %v", code))
				return
			}
			room.Join(player)
		} else {
			player.Room.Rejoin(player)
			room = player.Room
			str += " rejoined"
		}
		log.Println(str)
	default:
		log.Println("Invalid/missing action:", action)
		arcade.SendMsg(ws, "Invalid action: "+action)
		return
	}

	for {
		var msg arcade.PlayerMessage
		if err := ws.ReadJSON(&msg); err != nil {
			return
		}
		switch msg.Kind {
		case "end":
			room.EndTurn(player)
		case "seed":
			room.CastSeed(player, msg.Source, msg.Target)
		case "grow":
			room.GrowTree(player, msg.Target)
		case "sell":
			room.SellTree(player, msg.Target)
		}
	}
}
