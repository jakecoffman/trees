package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/trees/server/handlers"
	"github.com/jakecoffman/trees/server/lib"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO
		return true
	},
}

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, map[string]string{"Hello": "world"})
	})
	r.GET("/login", handlers.Login)
	r.GET("/ws", func(c *gin.Context) {
		playerId, err := c.Cookie("player")
		if err != nil {
			log.Println("Player failed to connect:", err.Error())
			c.String(401, "application/json", `{"error":"not logged in"}`)
			return
		}
		log.Println(playerId, "Connected")

		action := c.Query("action")
		code := c.Query("code")
		if action == "" {
			c.String(401, "application/json", `{"error":"missing action query string"}`)
			return
		}

		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("upgrade error:", err)
			return
		}
		defer ws.Close()

		handlers.WsHandler(lib.NewSafetySocket(ws), playerId, action, code)
	})
	log.Println("Serving http://127.0.0.1:8454")
	if err := http.ListenAndServe("127.0.0.1:8454", r); err != nil {
		log.Println(err)
	}
}
