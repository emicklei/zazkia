package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// Route represents the specification of a listener and client for connections (link).
type Route struct {
	Label           string          `json:"label"`
	ServiceHostname string          `json:"service-hostname"`
	ServicePort     int             `json:"service-port"`
	ServiceFromFile string          `json:"service-from-file"`
	ListenPort      int             `json:"listen-port"`
	Transport       *TransportState `json:"transport"`
}

// TransportState specifies how the transport of data is affected.
type TransportState struct {
	Verbose                      bool   `json:"verbose"`
	ThrottleServiceResponse      int    `json:"throttle-service-response"` // bytes per second
	DelayServiceResponse         int    `json:"delay-service-response"`    // milliseconds
	SendingToClient              bool   `json:"sending-to-client"`
	ReceivingFromClient          bool   `json:"receiving-from-client"`
	SendingToService             bool   `json:"sending-to-service"`
	ReceivingFromService         bool   `json:"receiving-from-service"`
	ServiceResponseCorruptMethod string `json:"service-response-corrupt-method"`
	AcceptConnections            bool   `json:"accept-connections"`
}

func readRoutes(location string) (routes []*Route, err error) {
	log.Println("reading routes from", location)
	f, err := os.Open(location)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&routes)
	// make sure transport is initialized
	for _, each := range routes {
		if each.Transport == nil {
			each.Transport = newDefaultTransportState()
		}
	}
	return
}

func (r Route) tcp() string {
	return fmt.Sprintf("%s:%d", r.ServiceHostname, r.ServicePort)
}

func (r Route) String() string {
	return fmt.Sprintf("[%s] :%d <-> %s:%d {%v}", r.Label, r.ListenPort, r.ServiceHostname, r.ServicePort, r.Transport)
}

func (s TransportState) String() string {
	return fmt.Sprintf("delay=%d,s2c=%v,rfc=%v,s2s=%v,rfs=%v,thr=%d,cor=%s,v=%v",
		s.DelayServiceResponse, s.SendingToClient, s.ReceivingFromClient, s.SendingToService, s.ReceivingFromService,
		s.ThrottleServiceResponse, s.ServiceResponseCorruptMethod, s.Verbose)
}

func newDefaultTransportState() *TransportState {
	return &TransportState{
		SendingToClient:      true,
		ReceivingFromClient:  true,
		SendingToService:     true,
		ReceivingFromService: true,
		AcceptConnections:    true,
	}
}
