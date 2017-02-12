package main

import (
	"strings"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
)

// RouteWithTags set the OpenAPI tags for this route.
func RouteWithTags(ws *restful.WebService, rb *restful.RouteBuilder) {
	rb.Metadata(restfulspec.KeyOpenAPITags, []string{strings.TrimLeft(ws.RootPath(), "/")})
	ws.Route(rb)
}
