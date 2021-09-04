package routes

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/crud"
	"github.com/jakecoffman/trees/server/arcade"
	"github.com/jakecoffman/trees/server/lib"
	"log"
	"net/http"
	"runtime/debug"
)

var Ws = crud.Spec{
	Method:      "GET",
	Path:        "/ws",
	Handler:     wsPreHandler,
	Description: "",
	Tags:        []string{"Play"},
	Summary:     "Connect to a game",
	Validate: crud.Validate{
		Query: crud.Object(map[string]crud.Field{
			"action": crud.String().Enum("new", "join").Description("Create a new game or join existing one"),
			"code":   crud.String().Min(6).Max(6).Description("The game room code"),
		}),
	},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO
		return true
	},
}

func wsPreHandler(c *gin.Context) {
	playerId, err := c.Cookie("player")
	if err != nil {
		log.Println("Player failed to connect:", err.Error())
		c.String(401, "application/json", `{"error":"not logged in"}`)
		c.Abort()
		return
	}
	log.Println(playerId, "Connected")

	action := c.Query("action")
	code := c.Query("code")
	if action == "" {
		c.String(401, "application/json", `{"error":"missing action query string"}`)
		c.Abort()
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		c.Abort()
		return
	}
	defer ws.Close()

	wsHandler(lib.NewSafetySocket(ws), playerId, action, code)
}

func wsHandler(ws *lib.SafetySocket, playerId, action, code string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Handler crashed:", r, string(debug.Stack()))
		} else {
			log.Println("Player disconnected")
		}
	}()

	player := arcade.Building.Enter(playerId, ws)
	defer arcade.Building.Disconnect(player)

	var room *arcade.Room

	// TODO move these into REST
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
