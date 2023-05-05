package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/irdaislakhuafa/pasino-api-forwarder/redirect"
)

type Handler interface {
	StartHandling(ctx context.Context)
}

type handler struct {
	server   *http.ServeMux
	redirect redirect.Redirect
}

func Init(ctx context.Context, server *http.ServeMux, redirect redirect.Redirect) Handler {
	return &handler{
		server:   server,
		redirect: redirect,
	}
}

func (handler *handler) StartHandling(ctx context.Context) {
	log.Println("Starting handle route...")

	handler.server.HandleFunc("/api/register", func(w http.ResponseWriter, r *http.Request) {
		handler.redirect.Redirect(ctx, w, r, "/api/register", nil)
	})

	handler.server.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		handler.redirect.Redirect(ctx, w, r, "/api/login", nil)
	})

	handler.server.HandleFunc("/account/get-socket-token", func(w http.ResponseWriter, r *http.Request) {
		handler.redirect.Redirect(ctx, w, r, "/account/get-socket-token", nil)
	})

	handler.server.HandleFunc("/deposit/get-deposit-information", func(w http.ResponseWriter, r *http.Request) {
		handler.redirect.Redirect(ctx, w, r, "/deposit/get-deposit-information", nil)
	})

	handler.server.HandleFunc("/withdraw/place-withdrawal", func(w http.ResponseWriter, r *http.Request) {
		handler.redirect.Redirect(ctx, w, r, "/withdraw/place-withdrawal", nil)
	})

	handler.server.HandleFunc("/transfer/send-transfer", func(w http.ResponseWriter, r *http.Request) {
		handler.redirect.Redirect(ctx, w, r, "/transfer/send-transfer", nil)
	})
	// TODO: redirect web socket
}
