package handlers

import (
	"log"
	"net/http"
)

const (
	PingEndpoint = "/ping"
	response     = "pong"
)

type PingHandler struct {
}

func (p *PingHandler) Get(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(response))
	if err != nil {
		log.Fatal("Error writing response in PingHandler")
	}
}
