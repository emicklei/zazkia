package main

import (
	"encoding/json"
	"fmt"
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
	Verbose                       bool     `json:"verbose"`
	ThrottleServiceResponse       int      `json:"throttle-service-response"` // bytes per second
	DelayServiceResponse          int      `json:"delay-service-response"`    // milliseconds
	SendingToClient               bool     `json:"sending-to-client"`
	ReceivingFromClient           bool     `json:"receiving-from-client"`
	SendingToService              bool     `json:"sending-to-service"`
	ReceivingFromService          bool     `json:"receiving-from-service"`
	ServiceResponseCorruptMethods []string `json:"service-response-corrupt-methods"`
}

func readRoutes() (routes []Route, err error) {
	f, err := os.Open("routes.json")
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
	return fmt.Sprintf("delay=%d,s2c=%v,rfc=%v,s2s=%v,rfs=%v",
		s.DelayServiceResponse, s.SendingToClient, s.ReceivingFromClient, s.SendingToService, s.ReceivingFromService)
}
