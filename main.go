package main

import (
	"log"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	server := http.DefaultServeMux

	log.Printf("Server listening on port %+v...\n", port)
	var handler http.Handler = server
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		panic(err)
	}
}
