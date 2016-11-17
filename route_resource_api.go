package main

import (
	"log"
	"net/http"
	"regexp"
)

var routesMatcher = regexp.MustCompile("/routes/(.+)/(.+)")

func (rr routeResource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/routes" {
		rr.routes(w, r)
		return
	}
	tokens := routesMatcher.FindStringSubmatch(r.URL.Path)
	if len(tokens) < 3 {
		http.NotFound(w, r)
		return
	}
	id := tokens[1]
	switch tokens[2] {
	case "toggle-accept":
		rr.toggleAcceptConnections(id, w, r)
	default:
		log.Println("unknown command", tokens[2])
		goHome(w, r)
	}
}
