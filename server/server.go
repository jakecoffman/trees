package main

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jakecoffman/trees/server/server"
	"github.com/jakecoffman/trees/server/server/lib"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Hello": "world!}"`))
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("player")
		if err == http.ErrNoCookie {
			playerId := uuid.New().String()
			cookie = &http.Cookie{
				Name:   "player",
				Value:  playerId,
				Path:   "/",
				MaxAge: 60,
			}
			http.SetCookie(w, cookie)
		}

		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		//w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Write([]byte(`{"status": "ok"}`))
	})
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer ws.Close()

		server.Handle(lib.NewSafetySocket(ws), r)
	})
	log.Println("Serving http://127.0.0.1:8080")
	if err := http.ListenAndServe("127.0.0.1:8080", mux); err != nil {
		log.Println(err)
	}
}
