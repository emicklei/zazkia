package main

//go:generate go-bindata -pkg main dashboard/...

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/elazarl/go-bindata-assetfs"
	restful "github.com/emicklei/go-restful"
)

var (
	oVerbose    = flag.Bool("v", false, "verbose logging")
	oAdminPort  = flag.Int("p", 9191, "port on which the admin http server will listen")
	oConfigfile = flag.String("f", "zazkia-routes.json", "route definition")

	routeMgr routeManager
	linkMgr  = newLinkManager()
)

func main() {
	log.Println("zazkia - tpc proxy for simulating network problems")
	flag.Parse()

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

}
