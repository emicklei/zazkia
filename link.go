package main

import (
	"fmt"
	"net"
)

// channel for consecutive integers to assign to new connections
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
	stats       TransportStats
}

type TransportStats struct {
	BytesSentToService       int64 `json:"bytesSentToService"`
	BytesSentToClient        int64 `json:"bytesSentToClient"`
	BytesReceivedFromClient  int64 `json:"bytesReceivedFromClient"`
	BytesReceivedFromService int64 `json:"bytesReceivedFromService"`
}

func (t TransportStats) String() string {
	return fmt.Sprintf("client sent:%d,service recv:%d,service sent:%d,client recv:%d",
		t.BytesReceivedFromClient,
		t.BytesSentToService,
		t.BytesReceivedFromService,
		t.BytesSentToClient)
}

func newLink(r Route, connectionToClient net.Conn, connectionToService net.Conn) *link {
	l := &link{
		ID:          <-idGen,
		route:       r,
		clientConn:  connectionToClient,
		serviceConn: connectionToService,
		stats:       TransportStats{},
	}
	l.resetTransport()
	// take the config from the route if given
	if r.hasTransportState() {
		l.transport = *r.Transport
	}
	return l
}

func (l *link) resetTransport() {
	l.transport = TransportState{
		Verbose:                      false,
		SendingToClient:              true,
		ReceivingFromClient:          true,
		SendingToService:             true,
		ReceivingFromService:         true,
		DelayServiceResponse:         0,
		ThrottleServiceResponse:      0,
		ServiceResponseCorruptMethod: "",
	}
}

func (l link) String() string {
	return fmt.Sprintf("[%s] %d: %s (s=%v,r=%v) <-> %s (s=%v,r=%v,d=%d,v=%v) [%s]",
		l.route.Label, l.ID,
		l.clientConn.RemoteAddr().String(), l.transport.SendingToService, l.transport.ReceivingFromService,
		l.serviceConn.RemoteAddr().String(), l.transport.SendingToClient, l.transport.ReceivingFromClient,
		l.transport.DelayServiceResponse, l.transport.Verbose,
		l.stats.String())
}

func (l *link) disconnect() error {
	if l.clientConn != nil {
		return l.clientConn.Close()
	}
	if l.serviceConn != nil {
		return l.serviceConn.Close()
	}
	return nil
}

type APILink struct {
	ID    int            `json:"id"`
	State TransportState `json:"state"`
	Stats TransportStats `json:"stats"`
}

func NewAPILink(l *link) APILink {
	return APILink{
		ID:    l.ID,
		State: l.transport,
		Stats: l.stats,
	}
}
