package main

import "net/http"

type linkResource struct {
	manager *linkManager
}

func (l linkResource) showLinks(id int, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	adminPage.Execute(w, l.manager.sortedLinks())
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
	gotoLinks(w, r)
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
		link.transport.DelayServiceResponse = 10000 // 10 sec
	}
	gotoLinks(w, r)
}

func (l linkResource) toggleSend(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.SendingToClient = !link.transport.SendingToClient
	link.transport.SendingToService = !link.transport.SendingToService
	gotoLinks(w, r)
}

func (l linkResource) toggleVerbose(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.Verbose = !link.transport.Verbose
	gotoLinks(w, r)
}

func (l linkResource) toggleReceive(id int, w http.ResponseWriter, r *http.Request) {
	link := getLink(id, w, r)
	if link == nil {
		return
	}
	link.transport.ReceivingFromClient = !link.transport.ReceivingFromClient
	link.transport.ReceivingFromService = !link.transport.ReceivingFromService
	gotoLinks(w, r)
}

func getLink(id int, w http.ResponseWriter, r *http.Request) *link {
	link := linkMgr.get(id)
	if link == nil {
		http.NotFound(w, r)
		return nil
	}
	return link
}
