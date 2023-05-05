package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/irdaislakhuafa/pasino-api-forwarder/handler"
	"github.com/irdaislakhuafa/pasino-api-forwarder/redirect"
)

var (
	port         = ""
	baseUrl      = "https://api.pasino.com"
	webSocketUrl = "wss://socket.pasino.com/dice/"
)

func main() {
	flag.StringVar(&port, "port", "8080", "Change default port (8080) of application")
	flag.Parse()

	ctx := context.Background()
	server := http.DefaultServeMux
	redirect := redirect.Init(ctx, baseUrl, http.DefaultClient, webSocketUrl)
	hand := handler.Init(ctx, server, redirect)

	hand.StartHandling(ctx)

	log.Printf("Server listening on port %+v...\n", port)
	var handler http.Handler = server
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		panic(err)
	}
}
