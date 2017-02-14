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
	"net"
)

// channel for consecutive integers
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
	clientConn   net.Conn
	serviceConn  net.Conn
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

func newLink(r *Route, connectionToClient net.Conn, connectionToService net.Conn) *link {
	l := &link{
		ID:          <-idGen,
		route:       r,
		clientConn:  connectionToClient,
		serviceConn: connectionToService,
	}
	if r.Transport == nil {
		l.resetTransport()
	} else {
		l.transport = *r.Transport
	}
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
	ID                  int            `json:"id"`
	State               TransportState `json:"state"`
	Stats               TransportStats `json:"stats"`
	ClientReceiveError  string         `json:"clientReceiveError,omitempty"`
	ServiceReceiveError string         `json:"serviceReceiveError,omitempty"`
}

func NewAPILink(l *link) APILink {
	return APILink{
		ID:                  l.ID,
		State:               l.transport,
		Stats:               l.stats,
		ClientReceiveError:  l.clientErrorString(),
		ServiceReceiveError: l.serviceErrorString(),
	}
}
