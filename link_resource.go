/*
Copyright 2017 Ernest Micklei

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful/v3"
)

type linkResource struct {
	manager *linkManager
}

func (l linkResource) addWebServiceTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/links")
	ws.Produces(restful.MIME_JSON)
	idParam := ws.PathParameter("id", "identifier of the link (is unique for all routes).")

	RouteWithTags(ws, ws.GET("/").To(l.getLinks).
		Doc("Returns all links for all routes."))
	RouteWithTags(ws, ws.POST("/closeAllWithError").To(l.closeLinksWithError).
		Doc("Cleanup all links for which an error was detected."))
	RouteWithTags(ws, ws.POST("/{id}/close").To(l.close).
		Doc("Close the client connections for this link. The server connection will be closed too.").
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/delay-response").To(l.delayResponse).
		Doc("Delay sending the response to the client").
		Param(ws.QueryParameter("ms", "milliseconds to delay transport to the client.").DataType("integer")).
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/break-response").To(l.breakResponse).
		Doc("Break sending the response to the client").
		Param(ws.QueryParameter("pct", "percentage of broken transport to the client.").DataType("integer")).
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/toggle-reads-client").To(l.toggleReadsClient).
		Doc("Change whether reading data from the client is enabled.").
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/toggle-writes-service").To(l.toggleWritesService).
		Doc("Change whether writing data to the service is enabled.").
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/toggle-reads-service").To(l.toggleReadsService).
		Doc("Change whether reading data from the service is enabled.").
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/toggle-writes-client").To(l.toggleWritesClient).
		Doc("Change whether writing data to the client is enabled.").
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/toggle-reads-client").To(l.toggleReadsClient).
		Doc("Change whether reading data from the client is enabled.").
		Param(idParam))
	RouteWithTags(ws, ws.POST("/{id}/toggle-verbose").To(l.toggleVerbose).
		Doc("Change whether transport of data (to and from the service) is logged.").
		Param(idParam))
	RouteWithTags(ws, ws.GET("/{id}/stats").To(l.stats).
		Doc("Statistics with data transported to and from services and clients.").
		Param(idParam))
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

func (l linkResource) breakResponse(request *restful.Request, response *restful.Response) {
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
	if link.transport.BreakServiceResponse > 0 {
		link.transport.BreakServiceResponse = 0
	} else {
		pct := request.QueryParameter("pct")
		msi, err := strconv.Atoi(pct)
		if msi < 0 {
			err = errors.New("pct parameter cannot be negative")
		}
		if msi > 100 {
			err = errors.New("pct parameter cannot be upper than 100")
		}
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		link.transport.BreakServiceResponse = msi
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
