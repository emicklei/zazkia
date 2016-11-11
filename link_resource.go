package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type linkResource struct {
	manager *linkManager
}

func (l linkResource) links(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(linkMgr.APIGroups())
}

func (l linkResource) closeLinksWithError(w http.ResponseWriter, r *http.Request) {
	l.manager.closeAllWithError()
	goHome(w, r)
}

func (l linkResource) close(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	if err := l.manager.disconnectAndRemove(link.ID); err != nil {
		http.NotFound(w, r)
		return
	}
	goHome(w, r)
}

func (l linkResource) delayResponse(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	// toggle
	if link.transport.DelayServiceResponse > 0 {
		link.transport.DelayServiceResponse = 0
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, err.Error())
			return
		}
		ms := r.Form.Get("ms")
		msi, err := strconv.Atoi(ms)
		if msi < 0 {
			err = errors.New("ms parameter cannot be negative")
		}
		if err != nil {
			log.Println(err.Error())
			w.WriteHeader(400)
			fmt.Fprintf(w, err.Error())
			return
		}
		link.transport.DelayServiceResponse = msi
	}
	goHome(w, r)
}

func (l linkResource) toggleReadsClient(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.ReceivingFromClient = !link.transport.ReceivingFromClient
	goHome(w, r)
}

func (l linkResource) toggleWritesService(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.SendingToService = !link.transport.SendingToService
	goHome(w, r)
}

func (l linkResource) toggleReadsService(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.ReceivingFromService = !link.transport.ReceivingFromService
	goHome(w, r)
}

func (l linkResource) toggleWritesClient(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.SendingToClient = !link.transport.SendingToClient
	goHome(w, r)
}

func (l linkResource) toggleVerbose(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.Verbose = !link.transport.Verbose
	goHome(w, r)
}

func (l linkResource) stats(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(link.stats)
}

func getLink(id int, w http.ResponseWriter, r *http.Request) *link {
	link := linkMgr.get(id)
	if link == nil {
		log.Println("no link with id", id)
		http.NotFound(w, r)
		return nil
	}
	return link
}
