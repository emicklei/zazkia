package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type routeResource struct {
	manager routeManager
}

func (rr routeResource) routes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rr.manager.routes)
}

func (rr routeResource) toggleAcceptConnections(id string, w http.ResponseWriter, r *http.Request) {
	var route *Route
	for _, each := range rr.manager.routes {
		if each.Label == id {
			route = each
			break
		}
	}
	if route == nil {
		http.NotFound(w, r)
		return
	}
	route.Transport.AcceptConnections = !route.Transport.AcceptConnections
	if route.Transport.AcceptConnections {
		log.Printf("start tcp listening for [%s]", route.Label)
	} else {
		log.Printf("stop tcp listening for [%s]", route.Label)
	}
	goHome(w, r)
}
