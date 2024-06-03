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
	"text/template"

	_ "embed"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/russross/blackfriday"
)

var indexPage *template.Template
var helpPage *template.Template

//go:embed dashboard/templates/index.html
var indexHtml string

//go:embed dashboard/templates/help.html
var helpHtml string

//go:embed dashboard/help.md
var helpMD []byte

func init() {
	indexPage = template.Must(template.New("index.html").Parse(indexHtml))
	helpPage = template.Must(template.New("help.html").Parse(helpHtml))
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
	if err := indexPage.Execute(response, linkMgr.APIGroups()); err != nil {
		log.Println("index rendering failed", err)
	}
}

func (d dashboardResource) getHelp(request *restful.Request, response *restful.Response) {
	response.Header().Set("Content-Type", "text/html")
	if err := helpPage.Execute(response, string(blackfriday.MarkdownCommon(helpMD))); err != nil {
		log.Println("help rendering failed", err)
	}
}
