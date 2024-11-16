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

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}

func (p *PingHandler) Get(writer http.ResponseWriter, req *http.Request) {
	_, err := writer.Write([]byte(response))
	if err != nil {
		log.Fatal("Error writing response in PingHandler")
	}
}
