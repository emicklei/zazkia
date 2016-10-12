package main

import (
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful"
)

type linkResource struct {
	manager *linkManager
}

func (r linkResource) showLinks(i *restful.Request, o *restful.Response) {
	o.ResponseWriter.Header().Set("Content-Type", "text/html")
	adminPage.Execute(o.ResponseWriter, r.manager.sortedLinks())
}

func (r linkResource) close(i *restful.Request, o *restful.Response) {
	link := getLink(i, o)
	if link == nil {
		return
	}
	if err := r.manager.disconnectAndRemove(link.ID); err != nil {
		o.WriteErrorString(http.StatusNotFound, err.Error())
		return
	}
	gotoLinks(i, o)
}

func (r linkResource) delayResponse(i *restful.Request, o *restful.Response) {
	link := getLink(i, o)
	if link == nil {
		return
	}
	// toggle
	if link.transport.DelayServiceResponse > 0 {
		link.transport.DelayServiceResponse = 0
	} else {
		link.transport.DelayServiceResponse = 10000 // 10 sec
	}
	gotoLinks(i, o)
}

func (r linkResource) toggleSend(i *restful.Request, o *restful.Response) {
	link := getLink(i, o)
	if link == nil {
		return
	}
	link.transport.SendingToClient = !link.transport.SendingToClient
	link.transport.SendingToService = !link.transport.SendingToService
	gotoLinks(i, o)
}

func (r linkResource) toggleReceive(i *restful.Request, o *restful.Response) {
	link := getLink(i, o)
	if link == nil {
		return
	}
	link.transport.ReceivingFromClient = !link.transport.ReceivingFromClient
	link.transport.ReceivingFromService = !link.transport.ReceivingFromService
	gotoLinks(i, o)
}

func gotoLinks(i *restful.Request, o *restful.Response) {
	http.Redirect(o.ResponseWriter, i.Request, "/links", http.StatusSeeOther)
}

func getLink(i *restful.Request, o *restful.Response) *link {
	id, err := strconv.Atoi(i.PathParameter("id"))
	if err != nil {
		o.WriteErrorString(http.StatusBadRequest, err.Error())
		return nil
	}
	link := linkMgr.get(id)
	if link == nil {
		o.WriteErrorString(http.StatusNotFound, "no such link")
		return nil
	}
	return link
}
