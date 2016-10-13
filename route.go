package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Route struct {
	Label         string
	ServiceDomain string          `json:"service-domain"`
	ServicePort   int             `json:"service-port"`
	ListenPort    int             `json:"listen-port"`
	Transport     *TransportState `json:"transport"`
}

type TransportState struct {
	Verbose                      bool   `json:"verbose"`
	ThrottleServiceResponse      int    `json:"throttle-service-response"` // bytes per second
	DelayServiceResponse         int    `json:"delay-service-response"`    // milliseconds
	SendingToClient              bool   `json:"sending-to-client"`
	ReceivingFromClient          bool   `json:"receiving-from-client"`
	SendingToService             bool   `json:"sending-to-service"`
	ReceivingFromService         bool   `json:"receiving-from-service"`
	ServiceResponseCorruptMethod string `json:"service-response-corrupt-method"`
}

func readRoutes(location string) (routes []Route, err error) {
	log.Println("reading routes from", location)
	f, err := os.Open(location)
	if err != nil {
		return
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&routes)
	return
}

func (r Route) tcp() string {
	return fmt.Sprintf("%s:%d", r.ServiceDomain, r.ServicePort)
}

func (r Route) hasTransportState() bool {
	return r.Transport != nil
}

func (r Route) String() string {
	return fmt.Sprintf("[%s] :%d <-> %s:%d {%v}", r.Label, r.ListenPort, r.ServiceDomain, r.ServicePort, r.Transport)
}

func (s TransportState) String() string {
	return fmt.Sprintf("delay=%d,s2c=%v,rfc=%v,s2s=%v,rfs=%v,thr=%d,cor=%s,v=%v",
		s.DelayServiceResponse, s.SendingToClient, s.ReceivingFromClient, s.SendingToService, s.ReceivingFromService,
		s.ThrottleServiceResponse, s.ServiceResponseCorruptMethod, s.Verbose)
}
