package main

import (
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/trees/server/server"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var upgrader = websocket.Upgrader{}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Hello": "world!}"`))
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer ws.Close()

		server.Handle(ws, r)
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
