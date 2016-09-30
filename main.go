package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var verbose = flag.Bool("v", false, "verbose logging")

var port = flag.Int("port", 49999, "port on which this service will listen")
var adminPort = flag.Int("admin.port", 48888, "port on which the admin http server will listen")

var routeMgr routeManager

var linkMgr = newLinkManager()

func main() {
	flag.Parse()

	// handle SIGINT (control+c)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		cleanAndExit(0)
	}()

	routes, err := readRoutes()
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

	log.Printf("start http listening on :%d\n", *adminPort)
	http.HandleFunc("/", commandHandler)
	err = http.ListenAndServe(fmt.Sprintf(":%d", *adminPort), nil)
	log.Println(err)
	cleanAndExit(1)
}

func cleanAndExit(code int) {
	log.Println("terminating...")
	linkMgr.close()
	os.Exit(code)
}

func acceptConnections(route Route, ln net.Listener) {
	log.Printf("start tcp listening for %#v", route)
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
		if err := transport(link, remoteConn, clientConn, !ReadsFromService); err != nil {
			log.Printf("failed to copy from client to remote:%v", err)
		}
		log.Printf("stopped writing to remote (%v), reading from client(%v)\n", addr, clientConn.RemoteAddr())
		linkMgr.disconnectAndRemove(link.ID)
	}()
	// client <- remote
	go func() {
		if err := transport(link, clientConn, remoteConn, ReadsFromService); err != nil {
			log.Printf("failed to copy from remote to client:%v", err)
		}
		log.Printf("stopped reading from remote (%v), writing to client (%v)\n", addr, clientConn.RemoteAddr())
		linkMgr.disconnectAndRemove(link.ID)
	}()
}

const ReadsFromService = true

var TransportBufferSize = 32 * 1024

func transport(link *link, w io.Writer, r io.Reader, readsFromService bool) error {
	// io.Copy with simulated problems
	buffer := make([]byte, TransportBufferSize)
	for {
		var (
			err  error
			read int
		)
		doRead := (readsFromService && link.receivingFromService) ||
			!readsFromService && link.receivingFromClient
		doWrite := (readsFromService && link.sendingToClient) ||
			!readsFromService && link.sendingToService

		if doRead {
			read, err = r.Read(buffer)
			if err != nil {
				return err
			}
		}
		if doWrite {
			offset := 0
			towrite := read
			for towrite > 0 {
				subset := buffer[offset:read]
				written, err := w.Write(subset)
				if err != nil {
					return err
				}
				if *verbose {
					log.Printf("[%s] written %d from %d", link.route.Label, written, read)
					log.Println(printable(subset))
				}
				offset += written
				towrite -= written
			}
		} else {
			if *verbose {
				log.Printf("[%s] flushing %d bytes", link.route.Label, read)
			}
		}
	}
}

func printable(data []byte) string {
	b := new(bytes.Buffer)
	for _, each := range data {
		if each == 10 || each == 13 { // CR,LF
			b.WriteByte(each)
			continue
		}
		if each >= 32 && each <= 126 {
			b.WriteByte(each)
		} else {
			b.WriteByte(46) // dot
		}
	}
	return b.String()
}
