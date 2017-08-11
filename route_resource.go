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
		Doc("All routes as defined in the configuration."))
	RouteWithTags(ws, ws.POST("/{label}/toggle-accept").To(rr.toggleAcceptConnections).
		Doc("Change whether new connections can be accepted for this route. Return").
		Param(labelParam))
	RouteWithTags(ws, ws.GET("/{label}/links").To(rr.getLinksForRoute).
		Doc("All links for connections created for this route.").
		Metadata("swagger.tags", "routes,links").
		Param(labelParam))
	RouteWithTags(ws, ws.POST("/{label}/links/close").To(rr.closeAllLinksForRoute).
		Doc("Close all connections created by links for this route.").
		Metadata("swagger.tags", "routes,links").
		Param(labelParam))
	container.Add(ws)
}

func (rr routeResource) getRoutes(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(rr.manager.routes)
}

func (rr routeResource) getLinksForRoute(request *restful.Request, response *restful.Response) {
	response.WriteAsJson(linkMgr.APILinks(request.PathParameter("label"))) // TODO should not use global linkMgr
}

func (rr routeResource) closeAllLinksForRoute(request *restful.Request, response *restful.Response) {
	label := request.PathParameter("label")
	log.Printf("[rest api] close all links for route [%s]", label)
	linkMgr.closeAllForRoute(label) // TODO should not use global linkMgr
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
	isAccepting := route.Transport.toggleAcceptConnections()
	if isAccepting {
		log.Printf("[rest api] start tcp listening for [%s]", route.Label)
	} else {
		log.Printf("[rest api] stop tcp listening for [%s]", route.Label)
	}
	response.WriteAsJson(isAccepting)
}
