package main

//go:generate go-bindata -pkg main dashboard/...

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"github.com/elazarl/go-bindata-assetfs"
)

var (
	oVerbose    = flag.Bool("v", false, "verbose logging")
	oAdminPort  = flag.Int("p", 9191, "port on which the admin http server will listen")
	oConfigfile = flag.String("f", "zazkia-routes.json", "route definition")

	routeMgr routeManager
	linkMgr  = newLinkManager()
)

func main() {
	flag.Parse()

	// handle SIGINT (control+c)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cleanAndExit(0)
	}()

	routes, err := readRoutes(*oConfigfile)
	if err != nil {
		log.Fatalf("failed to read routes:%v", err)
	}

	routeMgr = routeManager{routes: routes}

	// for each route start a listener
	for _, each := range routes {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", each.ListenPort))
		if err != nil {
			log.Fatalf("failed to start listener:%v", err)
		}
		go acceptConnections(each, ln)

	}

	log.Printf("start http listening on :%d\n", *oAdminPort)

	http.HandleFunc("/index.html", dashboardResource{}.index)
	dashboard := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: "dashboard"}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(dashboard)))
	http.Handle("/", linkResource{linkMgr})

	log.Println(http.ListenAndServe(fmt.Sprintf(":%d", *oAdminPort), nil))
	cleanAndExit(1)
}

func cleanAndExit(code int) {
	log.Println("terminating...")
	linkMgr.close()
	os.Exit(code)
}

func acceptConnections(route Route, ln net.Listener) {
	log.Printf("start tcp listening for %v", route)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept new connections:%v", err)
			break
		}
		go handleConnection(route, conn)
	}
}

func handleConnection(route Route, clientConn net.Conn) {
	addr := clientConn.RemoteAddr().String()
	serviceConn, err := net.Dial("tcp", route.tcp())
	if err != nil {
		log.Printf("[%s] failed to connect to remote:%v", route.Label, err)
		return
	}

	link := newLink(route, clientConn, serviceConn)
	linkMgr.add(link)

	log.Printf("[%s] start handling client(%v) <=> service(%v)\n", route.Label, addr, serviceConn.RemoteAddr())
	// service <- client
	go func() {
		if err := transport(link, serviceConn, clientConn, !AccessesService); err != nil {
			log.Printf("[%s] failed to copy from client to service:%v", route.Label, err)
		}
		log.Printf("[%s] stopped writing to service (%v), reading from client(%v)\n", route.Label, addr, clientConn.RemoteAddr())
	}()
	// client <- service
	go func() {
		if err := transport(link, clientConn, serviceConn, AccessesService); err != nil {
			log.Printf("[%s] failed to copy from service to client:%v", route.Label, err)
		}
		log.Printf("[%s] stopped reading from service (%v), writing to client (%v)\n", route.Label, addr, clientConn.RemoteAddr())
	}()
}
