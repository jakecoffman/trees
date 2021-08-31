package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"Hello": "world!}"`))
	})
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Println(err)
	}
}
