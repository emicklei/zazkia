package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
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
	remoteConn, err := net.Dial("tcp", route.tcp())
	if err != nil {
		log.Printf("[%s] failed to connect to remote:%v", route.Label, err)
		return
	}

	link := newLink(route, clientConn, remoteConn)
	linkMgr.add(link)

	log.Printf("start handling client(%v) <-> remote(%v)\n", addr, remoteConn.RemoteAddr())
	// remote <- client
	go func() {
		if err := transport(link, remoteConn, clientConn, !AccessesService); err != nil {
			log.Printf("failed to copy from client to remote:%v", err)
		}
		log.Printf("stopped writing to remote (%v), reading from client(%v)\n", addr, clientConn.RemoteAddr())
		linkMgr.disconnectAndRemove(link.ID)
	}()
	// client <- remote
	go func() {
		if err := transport(link, clientConn, remoteConn, AccessesService); err != nil {
			log.Printf("failed to copy from remote to client:%v", err)
		}
		log.Printf("stopped reading from remote (%v), writing to client (%v)\n", addr, clientConn.RemoteAddr())
		linkMgr.disconnectAndRemove(link.ID)
	}()
}
