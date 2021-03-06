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
	"strings"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
)

// RouteWithTags set the OpenAPI tags for this route.
func RouteWithTags(ws *restful.WebService, rb *restful.RouteBuilder) {
	rb.Metadata(restfulspec.KeyOpenAPITags, []string{strings.TrimLeft(ws.RootPath(), "/")})
	ws.Route(rb)
}
