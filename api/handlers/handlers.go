package handlers

import (
	"net/http"
)

const (
	PingEndpoint = "/ping"
)

func PingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
