package main

import (
	"fmt"
	"io"
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
	ID           int
	route        *Route
	clientConn   io.ReadWriteCloser
	clientLabel  string
	serviceConn  io.ReadWriteCloser
	serviceLabel string
	transport    TransportState
	stats        TransportStats
	clientError  error
	serviceError error
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

func newLink(r *Route, connectionToClient io.ReadWriteCloser, clientLabel string, connectionToService io.ReadWriteCloser, serviceLabel string) *link {
	l := &link{
		ID:          <-idGen,
		route:       r,
		clientConn:  connectionToClient,
		serviceConn: connectionToService,
		stats:       TransportStats{},
	}
	l.resetTransport()
	return l
}

func (l *link) clientErrorString() string {
	if l.clientError == nil {
		return ""
	}
	return l.clientError.Error()
}

func (l *link) serviceErrorString() string {
	if l.serviceError == nil {
		return ""
	}
	return l.serviceError.Error()
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
		l.clientLabel, l.transport.SendingToService, l.transport.ReceivingFromService,
		l.serviceLabel, l.transport.SendingToClient, l.transport.ReceivingFromClient,
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
	ID           int            `json:"id"`
	State        TransportState `json:"state"`
	Stats        TransportStats `json:"stats"`
	ClientError  string         `json:"clientError,omitempty"`
	ServiceError string         `json:"serviceError,omitempty"`
}

func NewAPILink(l *link) APILink {
	return APILink{
		ID:           l.ID,
		State:        l.transport,
		Stats:        l.stats,
		ClientError:  l.clientErrorString(),
		ServiceError: l.serviceErrorString(),
	}
}
