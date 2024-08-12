package router

import (
	"net/http"
)

type Route struct {
	handler  http.HandlerFunc
	method   string
	endpoint string
}

type router struct {
	http.Handler
	routes []Route
}

func NewRouter() *router {
	r := new(router)
	return r
}

func (r *router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var found int = -1

	for i, route := range r.routes {
		if route.endpoint == request.URL.Path {
			if request.Method == route.method {
				route.handler.ServeHTTP(writer, request)
				return
			} else {
				found = i
			}
		}
	}

	if found == -1 {
		writer.WriteHeader(http.StatusNotFound)
	} else {
		writer.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (r *router) AddRoute(handler http.HandlerFunc, method string, endpoint string) {
	r.routes = append(r.routes, Route{method: method, handler: handler, endpoint: endpoint})
}
