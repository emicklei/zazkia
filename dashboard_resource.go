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
	"text/template"

	restful "github.com/emicklei/go-restful"
	"github.com/russross/blackfriday"
)

var indexPage *template.Template
var helpPage *template.Template

func init() {
	t, err := Asset("dashboard/templates/index.html")
	if err != nil {
		panic("index.html")
	}
	indexPage = template.Must(template.New("index.html").Parse(string(t)))
	t, err = Asset("dashboard/templates/help.html")
	if err != nil {
		panic("help.html")
	}
	helpPage = template.Must(template.New("help.html").Parse(string(t)))
}

type dashboardResource struct{}

func (d dashboardResource) addWebServiceTo(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Produces(restful.MIME_JSON)
	ws.Route(ws.GET("/").To(d.getIndex))
	ws.Route(ws.GET("/index.html").To(d.getIndex))
	ws.Route(ws.GET("/help.html").To(d.getHelp))
	container.Add(ws)
}

func (d dashboardResource) getIndex(request *restful.Request, response *restful.Response) {
	response.Header().Set("Content-Type", "text/html")
	indexPage.Execute(response, linkMgr.APIGroups())
}

func (d dashboardResource) getHelp(request *restful.Request, response *restful.Response) {
	response.Header().Set("Content-Type", "text/html")
	input, err := Asset("dashboard/help.md")
	if err != nil {
		response.WriteHeaderAndEntity(404, "could not load/find help.md")
		return
	}
	helpPage.Execute(response, string(blackfriday.MarkdownCommon(input)))
}
