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
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func startListeners(routes []*Route) {
	// for each route start a listener
	for _, each := range routes {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", each.ListenPort))
		if err != nil {
			log.Fatalf("failed to start listener:%v", err)
		}
		go acceptConnections(each, ln)
	}
}

func acceptConnections(route *Route, ln net.Listener) {
	log.Printf("start tcp listening for %v", route)
	for {
		if !route.Transport.isAcceptConnections() {
			log.Printf("not accepting new connections for %s, retrying in 1 second", route.Label)
			time.Sleep(1 * time.Second)
			continue
		}
		// Accept is blocking so changing the transport will not prevent from getting a connection.
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("failed to accept new connections for %s because [%v]", route.Label, err)
			break
		}
		// Is Accept was disbled then close it
		if !route.Transport.isAcceptConnections() {
			log.Printf("got new connection but accept connections is not enabled for [%v], so I will close it immediately", route.Label)
			if err := conn.Close(); err != nil {
				log.Printf("failed to close connection for %s because [%v]", route.Label, err)
			}
			continue
		}
		go handleConnection(route, conn)
	}
}

func handleConnection(route *Route, clientConn net.Conn) {
	addr := clientConn.RemoteAddr().String()
	serviceConn, err := net.Dial("tcp", route.tcp())
	if err != nil {
		log.Printf("[%s] failed to connect to remote:%v", route.Label, err)
		return
	}

	link := newLink(route, clientConn, serviceConn)
	linkMgr.add(link)
	
	log.Printf("[%s:%d] start handling client(%v) <=> service(%v)\n", route.Label, link.ID, addr, serviceConn.RemoteAddr())
	// service <- client
	go func() {
		if err := transport(link, serviceConn, clientConn, !AccessesService); err != nil {
			if( err == io.EOF ) {
				log.Printf("[%s:%d] stopped writing to service (%v), end of reading from client(%v)\n", route.Label, link.ID, serviceConn.RemoteAddr(), clientConn.RemoteAddr())
			} else {
				log.Printf("[%s:%d] stopped writing to service (%v), reading from client(%v), with error (%v)\n", route.Label, link.ID, serviceConn.RemoteAddr(), clientConn.RemoteAddr(), err)
			}
			link.clientError = err
		}
		if( *oClose ) {
			time.Sleep(time.Duration(1000) * time.Millisecond)
			linkMgr.disconnectAndRemove(link.ID)
		}
	}()
	// client <- service
	go func() {
		if err := transport(link, clientConn, serviceConn, AccessesService); err != nil {
			if( err == io.EOF ) {
				log.Printf("[%s:%d] stopped reading from service (%v), end of writing to client (%v)\n", route.Label, link.ID, serviceConn.RemoteAddr(), clientConn.RemoteAddr())
			} else {
				log.Printf("[%s:%d] stopped reading from service (%v), writing to client (%v), with error (%v)\n", route.Label, link.ID, serviceConn.RemoteAddr(), clientConn.RemoteAddr(), err)
			}
			link.serviceError = err
		}
		if( *oClose ) {
			linkMgr.disconnectAndRemove(link.ID)
		}
	}()
}
