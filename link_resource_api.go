package main

import "github.com/emicklei/go-restful"

var noOperation = func(*restful.Request, *restful.Response) {}

func registerLinkResource(l linkResource, container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		ApiVersion("1.0").
		Path("/v1/links").
		Doc("Link management").
		Consumes("*/*").
		Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/").To(l.showLinks))
	ws.Route(ws.GET("/{id}/close").To(l.close))
	ws.Route(ws.GET("/{id}/delay-response").To(l.delayResponse))
	ws.Route(ws.GET("/{id}/toggle-send").To(l.toggleSend))
	ws.Route(ws.GET("/{id}/toggle-receive").To(l.toggleReceive))
	ws.Route(ws.GET("/favion.ico").To(noOperation))
}
