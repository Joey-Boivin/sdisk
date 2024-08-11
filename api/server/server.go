package server

import (
	"net/http"

	"github.com/Joey-Boivin/cdisk/api/handlers"
)

type server struct {
	http.Handler
}

func NewServer() *server {
	mux := http.NewServeMux()
	mux.Handle(handlers.PingEndpoint, http.HandlerFunc(handlers.PingHandler))
	r := new(server)
	r.Handler = mux
	return r
}
