package main

import (
	"github.com/gin-gonic/gin"
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
	gin.SetMode(gin.ReleaseMode)
}

func main() {
	engine := gin.Default()
	a := &adapter.Adapter{Engine: engine}
	r := crud.NewRouter("Trees", "1.0.0", a)
	r.Swagger.BasePath = "/api"

	if err := r.Add(routes.Login); err != nil {
		log.Fatal(err)
	}
	if err := r.Add(routes.Ws); err != nil {
		log.Fatal(err)
	}
	if err := r.Add(routes.Room...); err != nil {
		log.Fatal(err)
	}
	if err := r.Add(routes.Admin); err != nil {
		log.Fatal(err)
	}

	log.Println("Serving http://127.0.0.1:8454")
	if err := r.Serve("127.0.0.1:8454"); err != nil {
		log.Println(err)
	}
}
