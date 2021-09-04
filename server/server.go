package main

import (
	"github.com/jakecoffman/crud"
	adapter "github.com/jakecoffman/crud/adapters/gin-adapter"
	"github.com/jakecoffman/trees/server/routes"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	r := crud.NewRouter("Trees", "1.0.0", adapter.New())

	if err := r.Add(routes.Login); err != nil {
		log.Fatal(err)
	}
	if err := r.Add(routes.Ws); err != nil {
		log.Fatal(err)
	}

	log.Println("Serving http://127.0.0.1:8454")
	if err := r.Serve("127.0.0.1:8454"); err != nil {
		log.Println(err)
	}
}
