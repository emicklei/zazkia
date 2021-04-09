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

//go:generate go-bindata -pkg main dashboard/... swagger-ui/...

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
)

var (
	oClose      = flag.Bool("c", false, "automatic connection close")
	oVerbose    = flag.Bool("v", false, "verbose logging")
	oAdminPort  = flag.Int("p", 9191, "port on which the admin http server will listen")
	oConfigfile = flag.String("f", "zazkia-routes.json", "route definition")

	routeMgr routeManager
	linkMgr  = newLinkManager()
)

func main() {
	log.Println("zazkia - tpc proxy for simulating network problems")
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	// handle SIGINT (control+c)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cleanAndExit(0)
	}()

	// reading config
	routes, err := readRoutes(*oConfigfile)
	if err != nil {
		here, _ := os.Getwd()
		log.Fatalf("failed to read routes in %s because [%v]", here, err)
	} else {
		log.Println("done reading routes from", *oConfigfile)
	}

	routeMgr = routeManager{routes: routes}

	// for each route start a listener
	startListeners(routes)

	log.Printf("start http listening on :%d\n", *oAdminPort)

	// static file serving
	dashboard := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "dashboard"}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(dashboard)))

	addRESTResources()
	addSwagger()

	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *oAdminPort), nil))
	cleanAndExit(1)
}

func addRESTResources() {
	restful.DefaultRequestContentType(restful.MIME_JSON)
	restful.DefaultResponseContentType(restful.MIME_JSON)
	routeResource{routeMgr}.addWebServiceTo(restful.DefaultContainer)
	dashboardResource{}.addWebServiceTo(restful.DefaultContainer)
	linkResource{linkMgr}.addWebServiceTo(restful.DefaultContainer)
}

func cleanAndExit(code int) {
	log.Println("terminating...")
	linkMgr.close()
	os.Exit(code)
}

func addSwagger() {
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		WebServicesURL:                "http://localhost:" + strconv.Itoa(*oAdminPort),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: extendSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// static file serving
	swaggerUI := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "swagger-ui/dist"}
	http.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(swaggerUI)))
}

func extendSwaggerObject(s *spec.Swagger) {
	s.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Zazkia REST API",
			Description: "Resources for managing Routes and Links",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name: "Ernest Micklei",
					URL:  "https://github.com/emicklei/zazkia",
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: "MIT License",
					URL:  "https://opensource.org/licenses/MIT",
				},
			},
			Version: "1.0.0",
		},
	}
	s.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "links",
		Description: "Managing links"}}, {TagProps: spec.TagProps{
		Name:        "routes",
		Description: "Managing routes"}}}
}
