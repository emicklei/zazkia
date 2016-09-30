package main

import (
	"fmt"
	"net"
)

var idGen chan int

func init() {
	idGen = make(chan int)
	go func() {
		id := 1
		for {
			idGen <- id
			id++
		}
	}()
}

type link struct {
	ID                   int
	route                Route
	clientConn           net.Conn
	serviceConn          net.Conn
	sendingToClient      bool
	receivingFromClient  bool
	sendingToService     bool
	receivingFromService bool
}

func newLink(r Route, connectionToClient net.Conn, connectionToService net.Conn) *link {
	return &link{
		ID:                   <-idGen,
		route:                r,
		clientConn:           connectionToClient,
		serviceConn:          connectionToService,
		sendingToClient:      true,
		receivingFromClient:  true,
		sendingToService:     true,
		receivingFromService: true,
	}
}

func (l link) String() string {
	return fmt.Sprintf("[%s] %d: %s (s=%v,r=%v) <-> %s (s=%v,r=%v)",
		l.route.Label, l.ID,
		l.clientConn.RemoteAddr().String(), l.sendingToService, l.receivingFromService,
		l.serviceConn.RemoteAddr().String(), l.sendingToClient, l.receivingFromClient)
}

func (l *link) disconnect() {
	if l.clientConn != nil {
		l.clientConn.Close()
	}
	if l.serviceConn != nil {
		l.serviceConn.Close()
	}
}
