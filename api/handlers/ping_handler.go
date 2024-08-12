package handlers

import (
	"net/http"
)

const (
	PingEndpoint = "/ping"
	response     = "pong"
)

type PingHandler struct {
}

func (p *PingHandler) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(response))
}
