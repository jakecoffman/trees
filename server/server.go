package main

import (
	"log"
	"net/http"
)

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
