package main

import (
	"log"
	"net/http"

	restful "github.com/emicklei/go-restful"
)

type routeResource struct {
	manager routeManager
}

func (rr routeResource) addWebServiceTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/routes")
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(rr.getRoutes))
	ws.Route(ws.POST("/{id}/toggle-accept").To(rr.toggleAcceptConnections))
	container.Add(ws)
}

func (rr routeResource) getRoutes(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(rr.manager.routes)
}

func (rr routeResource) toggleAcceptConnections(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("id")
	var route *Route
	for _, each := range rr.manager.routes {
		if each.Label == id {
			route = each
			break
		}
	}
	if route == nil {
		http.NotFound(response, request.Request)
		return
	}
	route.Transport.AcceptConnections = !route.Transport.AcceptConnections
	if route.Transport.AcceptConnections {
		log.Printf("start tcp listening for [%s]", route.Label)
	} else {
		log.Printf("stop tcp listening for [%s]", route.Label)
	}
}
