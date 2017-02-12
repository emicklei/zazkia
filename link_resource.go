package main

import (
	"errors"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
)

type linkResource struct {
	manager *linkManager
}

func (l linkResource) addWebServiceTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/links")
	ws.Produces(restful.MIME_JSON)
	RouteWithTags(ws, ws.GET("/").To(l.getLinks))
	RouteWithTags(ws, ws.POST("/closeAllWithError").To(l.closeLinksWithError))
	RouteWithTags(ws, ws.POST("/{id}/close").To(l.close))
	RouteWithTags(ws, ws.POST("/{id}/delay-response").To(l.delayResponse))
	RouteWithTags(ws, ws.POST("/{id}/toggle-reads-client").To(l.toggleReadsClient))
	RouteWithTags(ws, ws.POST("/{id}/toggle-writes-service").To(l.toggleWritesService))
	RouteWithTags(ws, ws.POST("/{id}/toggle-reads-service").To(l.toggleReadsService))
	RouteWithTags(ws, ws.POST("/{id}/toggle-writes-client").To(l.toggleWritesClient))
	RouteWithTags(ws, ws.POST("/{id}/toggle-reads-client").To(l.toggleReadsClient))
	RouteWithTags(ws, ws.POST("/{id}/toggle-verbose").To(l.toggleVerbose))
	RouteWithTags(ws, ws.GET("/{id}/stats").To(l.stats))
	container.Add(ws)
}

func (l linkResource) getLinks(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(l.manager.APIGroups())
}

func (l linkResource) closeLinksWithError(request *restful.Request, response *restful.Response) {
	l.manager.closeAllWithError()

}

func (l linkResource) close(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	if err := l.manager.disconnectAndRemove(link.ID); err != nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}

}

func (l linkResource) delayResponse(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	// toggle
	if link.transport.DelayServiceResponse > 0 {
		link.transport.DelayServiceResponse = 0
	} else {
		ms := request.QueryParameter("ms")
		msi, err := strconv.Atoi(ms)
		if msi < 0 {
			err = errors.New("ms parameter cannot be negative")
		}
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		link.transport.DelayServiceResponse = msi
	}

}

func (l linkResource) toggleReadsClient(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	link.transport.ReceivingFromClient = !link.transport.ReceivingFromClient

}

func (l linkResource) toggleWritesService(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	link.transport.SendingToService = !link.transport.SendingToService

}

func (l linkResource) toggleReadsService(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	link.transport.ReceivingFromService = !link.transport.ReceivingFromService

}

func (l linkResource) toggleWritesClient(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	link.transport.SendingToClient = !link.transport.SendingToClient

}

func (l linkResource) toggleVerbose(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	link.transport.Verbose = !link.transport.Verbose

}

func (l linkResource) stats(request *restful.Request, response *restful.Response) {
	id, ok := linkIDFromRequest(request, response)
	if !ok {
		return
	}
	link := l.manager.get(id)
	if link == nil {
		response.WriteHeader(http.StatusNotFound)
		return
	}
	response.WriteAsJson(link.stats)
}

func linkIDFromRequest(request *restful.Request, response *restful.Response) (int, bool) {
	ids := request.PathParameter("id")
	id, err := strconv.Atoi(ids)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
