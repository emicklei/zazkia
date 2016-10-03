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
	ID          int
	route       Route
	clientConn  net.Conn
	serviceConn net.Conn
	transport   TransportState
}

func newLink(r Route, connectionToClient net.Conn, connectionToService net.Conn) *link {
	l := &link{
		ID:          <-idGen,
		route:       r,
		clientConn:  connectionToClient,
		serviceConn: connectionToService,
		transport: TransportState{
			SendingToClient:      true,
			ReceivingFromClient:  true,
			SendingToService:     true,
			ReceivingFromService: true,
			DelayServiceResponse: 0,
		},
	}
	if r.hasTransportState() {
		l.transport = *r.Transport
	}
	return l
}

func (l link) String() string {
	return fmt.Sprintf("[%s] %d: %s (s=%v,r=%v) <-> %s (s=%v,r=%v)",
		l.route.Label, l.ID,
		l.clientConn.RemoteAddr().String(), l.transport.SendingToService, l.transport.ReceivingFromService,
		l.serviceConn.RemoteAddr().String(), l.transport.SendingToClient, l.transport.ReceivingFromClient)
}

func (l *link) disconnect() {
	if l.clientConn != nil {
		l.clientConn.Close()
	}
	if l.serviceConn != nil {
		l.serviceConn.Close()
	}
}
