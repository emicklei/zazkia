package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
)

var linksMatcher = regexp.MustCompile("/links/(\\d*)/(.+)")

func (l linkResource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/index.html", http.StatusMovedPermanently)
		return
	}
	if r.URL.Path == "/links" {
		l.links(w, r)
		return
	}
	if r.URL.Path == "/links/closeAllWithError" {
		l.closeLinksWithError(w, r)
		return
	}
	tokens := linksMatcher.FindStringSubmatch(r.URL.Path)
	if len(tokens) < 3 {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(tokens[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch tokens[2] {
	case "close":
		l.close(id, w, r)
	case "delay-response":
		l.delayResponse(id, w, r)
	case "toggle-reads-client":
		l.toggleReadsClient(id, w, r)
	case "toggle-reads-service":
		l.toggleReadsService(id, w, r)
	case "toggle-writes-client":
		l.toggleWritesClient(id, w, r)
	case "toggle-writes-service":
		l.toggleWritesService(id, w, r)
	case "toggle-verbose":
		l.toggleVerbose(id, w, r)
	case "stats":
		l.stats(id, w, r)
	default:
		log.Println("unknown command", tokens[2])
		goHome(w, r)
	}
}

func goHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
