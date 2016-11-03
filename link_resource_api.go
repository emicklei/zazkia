package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
)

var linksMatcher = regexp.MustCompile("/links/(\\d*)/(.+)")

func (l linkResource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.Header().Set("Content-Type", "text/html")
		adminPage.Execute(w, linkMgr.APIGroups())
		return
	}
	if r.URL.Path == "/links" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(linkMgr.APIGroups())
		return
	}
	tokens := linksMatcher.FindStringSubmatch(r.URL.Path)
	if len(tokens) != 3 {
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
	case "toggle-receive":
		l.toggleReceive(id, w, r)
	case "toggle-send":
		l.toggleSend(id, w, r)
	case "toggle-verbose":
		l.toggleVerbose(id, w, r)
	default:
		goHome(w, r)
	}
}

func goHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
