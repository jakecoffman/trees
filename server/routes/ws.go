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
	"sync/atomic"
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

	code := c.Query("code")

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		c.Abort()
		return
	}
	defer ws.Close()

	wsHandler(lib.NewSafetySocket(ws), playerId, code)
}

func wsHandler(ws *lib.SafetySocket, playerId, code string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Handler crashed:", r, string(debug.Stack()))
		} else {
			log.Println("Player disconnected")
		}
	}()
	atomic.AddUint64(&arcade.ActiveWsConnections, 1)
	defer atomic.AddUint64(&arcade.ActiveWsConnections, -1)

	player := arcade.Building.Enter(playerId, ws)
	defer arcade.Building.Disconnect(player)

	var room *arcade.Room

	str := fmt.Sprint(playerId)
	// player has joined a different room, quit/forfeit current game
	if player.Room != nil && player.Room.Code != code {
		str += " leaving"
		player.Room.Quit(player)
	}
	// player is joining fresh
	if player.Room == nil {
		room = arcade.Building.GetRoom(code)
		str += " fresh join"
		if room == nil {
			arcade.SendMsg(ws, fmt.Sprintf("Code is wrong, or room is gone: %v", code))
			return
		}
		room.Join(player)
	} else {
		// player is rejoining
		player.Room.Rejoin(player)
		room = player.Room
		str += " rejoined"
	}
	log.Println(str)

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
		case "bot":
			room.UseBot()
		}
	}
}
