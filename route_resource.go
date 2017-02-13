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
	labelParam := ws.PathParameter("label", "identifier of the route.")

	RouteWithTags(ws, ws.GET("/").To(rr.getRoutes).
		Doc("All routes as defined in the configuration.").
		Param(labelParam))
	RouteWithTags(ws, ws.POST("/{label}/toggle-accept").To(rr.toggleAcceptConnections).
		Doc("Change whether new connections can be accepted for this route.").
		Param(labelParam))
	RouteWithTags(ws, ws.GET("/{label}/links").To(rr.getLinksForRoute).
		Doc("All links for connections created for this route.").
		Metadata("swagger.tags", "routes,links").
		Param(labelParam))
	container.Add(ws)
}

func (rr routeResource) getRoutes(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(rr.manager.routes)
}

func (rr routeResource) getLinksForRoute(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("label")
	response.WriteAsJson(linkMgr.APILinks(id)) // TODO should not use global
}

func (rr routeResource) toggleAcceptConnections(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("label")
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
